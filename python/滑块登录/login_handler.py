from __future__ import annotations

import io
import os
import random
import re
import time
from datetime import datetime
from typing import Callable, Tuple

from PIL import Image
import cv2
import numpy as np
from selenium import webdriver
from selenium.common.exceptions import TimeoutException
from selenium.webdriver import ActionChains
from selenium.webdriver.chrome.options import Options as ChromeOptions
from selenium.webdriver.common.by import By
from selenium.webdriver.edge.options import Options as EdgeOptions
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.ui import WebDriverWait

LogCallback = Callable[[str], None]


class SmartEduLogin:
    """SmartEdu 登录处理器，使用 OpenCV 进行滑块验证码识别。"""

    def __init__(self, log_callback: LogCallback | None = None, browser: str = "chrome") -> None:
        self.log_callback = log_callback
        self.browser = browser
        self.driver: webdriver.Remote | None = None
        self.login_url = "https://auth.smartedu.cn/uias/login"

    def log(self, msg: str) -> None:
        """输出日志信息。"""
        if self.log_callback:
            self.log_callback(msg)

    def login(self, username: str, password: str) -> Tuple[bool, str, float | None]:
        """执行完整登录流程，并获取用户积分。"""
        points = None
        try:
            # 创建浏览器驱动
            self.driver = self._create_driver(self.browser)
            self.log("启动浏览器，打开登录页面...")
            self.driver.get(self.login_url)

            wait = WebDriverWait(self.driver, 15)

            # 勾选用户协议
            try:
                agreement = wait.until(EC.element_to_be_clickable((By.ID, "agreementCheckbox")))
                if not agreement.is_selected():
                    agreement.click()
                    self.log("已勾选用户协议")
            except TimeoutException:
                self.log("未找到协议复选框，跳过")

            # 输入账号密码
            user_input = wait.until(EC.visibility_of_element_located((By.ID, "username")))
            pwd_input = wait.until(EC.visibility_of_element_located((By.ID, "tmpPassword")))
            user_input.clear()
            user_input.send_keys(username)
            pwd_input.clear()
            pwd_input.send_keys(password)
            self.log("已输入账号密码")

            # 点击登录按钮，触发滑块验证
            login_btn = wait.until(EC.element_to_be_clickable((By.ID, "loginBtn")))
            login_btn.click()
            self.log("已点击登录按钮，等待滑块验证码...")

            # 处理滑块验证码
            if not self._handle_slider_captcha():
                return False, "滑块验证失败或未通过", None

            # 等待页面跳转
            time.sleep(3)

            # 检查登录结果
            current_url = self.driver.current_url
            page_source = self.driver.page_source
            self.log(f"当前地址: {current_url}")

            if "basic.smartedu.cn" in current_url:
                # 登录成功，获取用户积分
                points = self.fetch_user_points()
                return True, "登录成功", points

            if "账号或密码" in page_source or "请输入账号或密码" in page_source:
                return False, "账号或密码错误", None

            return False, "登录失败，请检查账号密码或网络连接", None

        except Exception as exc:  # noqa: BLE001
            self.log(f"登录流程异常: {exc}")
            return False, str(exc), None
        finally:
            if self.driver is not None:
                self.driver.quit()
                self.driver = None

    def _create_driver(self, browser: str) -> webdriver.Remote:
        """创建浏览器驱动。"""
        if browser.lower() == "edge":
            options = EdgeOptions()
            options.add_argument("--start-maximized")
            return webdriver.Edge(options=options)

        # 默认使用 Chrome
        options = ChromeOptions()
        options.add_argument("--start-maximized")
        return webdriver.Chrome(options=options)

    def _handle_slider_captcha(self) -> bool:
        """处理滑块验证码。"""
        assert self.driver is not None
        wait = WebDriverWait(self.driver, 15)

        try:
            # 等待并切换到滑块验证码 iframe
            iframe = wait.until(
                EC.presence_of_element_located((By.CSS_SELECTOR, "iframe#tcaptcha_iframe_dy"))
            )
            self.log("检测到滑块验证码 iframe，准备切换...")
            self.driver.switch_to.frame(iframe)

            # 等待 iframe 内容完全加载：先等待背景图元素出现（作为加载完成的标志）
            try:
                iframe_wait = WebDriverWait(self.driver, 10)
                iframe_wait.until(EC.presence_of_element_located((By.ID, "slideBg")))
                self.log("iframe 内容已加载（背景图已出现）")
            except TimeoutException:
                self.log("等待背景图超时，继续尝试...")

            # 额外等待一段时间，确保所有元素（包括滑块）都已渲染
            time.sleep(2.0)
            self.log("等待 iframe 渲染完成...")

            # 查找滑块元素（带重试机制）
            slider = self._find_slider_element()
            if slider is None:
                # 如果第一次找不到，再等待更长时间后重试一次
                self.log("首次未找到滑块元素，等待 3 秒后重试...")
                time.sleep(3.0)
                slider = self._find_slider_element()
                
            if slider is None:
                self.log("重试后仍未找到滑块元素")
                self.driver.switch_to.default_content()
                return False

            self.log("已定位到滑块元素")

            # 定位背景图和拼图块并使用 OpenCV 识别缺口
            try:
                bg_elem = wait.until(EC.presence_of_element_located((By.ID, "slideBg")))

                # 查找拼图块（约 50x50px）
                all_items = self.driver.find_elements(By.CSS_SELECTOR, ".tc-fg-item")
                piece_elem = None
                for item in all_items:
                    classes = item.get_attribute("class") or ""
                    if "tc-slider-normal" not in classes:
                        size = item.size
                        if 45 <= size["width"] <= 55 and 45 <= size["height"] <= 55:
                            piece_elem = item
                            break

                if piece_elem is None:
                    raise RuntimeError("未找到拼图块元素")

                # 仅使用自写的 OpenCV 图像处理逻辑识别缺口
                distance = self._detect_gap_with_cv(bg_elem, piece_elem)
                self.log(f"OpenCV 识别到需要拖动距离: {distance}px")

            except Exception as exc:  # noqa: BLE001
                self.log(f"识别失败: {exc}")
                return False

            # 执行拖动
            self._drag_slider(slider, distance)

            # 等待验证结果
            self.driver.switch_to.default_content()
            time.sleep(2)
            return True

        except TimeoutException:
            self.log("滑块验证码元素定位超时")
            try:
                self.driver.switch_to.default_content()
            except Exception:  # noqa: BLE001
                pass
            return False
        except Exception as exc:  # noqa: BLE001
            self.log(f"滑块处理异常: {exc}")
            try:
                self.driver.switch_to.default_content()
            except Exception:  # noqa: BLE001
                pass
            return False

    def _find_slider_element(self):
        """查找滑块元素（优先按固定类名 `tc-fg-item tc-slider-normal` 查找）。"""
        assert self.driver is not None

        # 1）优先等待固定类名出现（增加等待时间到 10 秒）
        try:
            wait = WebDriverWait(self.driver, 10)
            elem = wait.until(
                EC.visibility_of_element_located((By.CSS_SELECTOR, ".tc-fg-item.tc-slider-normal"))
            )
            if elem.is_displayed():
                self.log("通过 WebDriverWait 找到滑块元素 (.tc-fg-item.tc-slider-normal)")
                return elem
        except Exception:  # noqa: BLE001
            # 不打印异常，进入自定义轮询
            pass

        # 2）自定义轮询：仍以固定类名为主，其它选择器仅作为兜底
        selectors = [
            ".tc-fg-item.tc-slider-normal",
            ".tc-slider-normal",  # 兜底：有时可能少一个类
        ]

        # 进一步轮询更长时间，兼容加载很慢的情况（增加到 60 次，约 6 秒）
        for i in range(60):  # 约 6 秒
            found = None
            used_selector = None
            for selector in selectors:
                try:
                    sliders = self.driver.find_elements(By.CSS_SELECTOR, selector)
                    for s in sliders:
                        if s.is_displayed():
                            found = s
                            used_selector = selector
                            break
                    if found is not None:
                        break
                except Exception:  # noqa: BLE001
                    continue

            if found is not None:
                self.log(f"通过选择器 {used_selector} 找到滑块元素（轮询第 {i+1} 次）")
                return found

            time.sleep(0.1)

        # 如果仍然未找到，输出调试信息，方便你查看页面结构
        try:
            all_items = self.driver.find_elements(By.CSS_SELECTOR, ".tc-fg-item")
            self.log(f"[调试] 轮询结束仍未找到滑块，一共找到 {len(all_items)} 个 .tc-fg-item：")
            for idx, item in enumerate(all_items[:5]):
                try:
                    classes = item.get_attribute("class")
                    size = item.size
                    displayed = item.is_displayed()
                    self.log(f"  fg-item[{idx}]: class={classes}, size={size}, displayed={displayed}")
                except Exception:  # noqa: BLE001
                    continue
        except Exception:  # noqa: BLE001
            pass

        return None

    def _detect_gap_with_cv(self, bg_elem, piece_elem) -> int:
        """使用自写 OpenCV 图像处理逻辑识别缺口位置，返回需要拖动的距离（DOM 坐标）。"""
        try:
            # 计算拼图块当前在 DOM 中相对于背景图的 X 位置
            bg_loc = bg_elem.location
            piece_loc = piece_elem.location
            current_x_dom = piece_loc["x"] - bg_loc["x"]

            # 截图背景图和拼图块
            bg_png = bg_elem.screenshot_as_png
            piece_png = piece_elem.screenshot_as_png

            bg_img = Image.open(io.BytesIO(bg_png))
            piece_img = Image.open(io.BytesIO(piece_png))

            # 转为 OpenCV 图像（灰度）
            bg_cv = cv2.cvtColor(np.array(bg_img), cv2.COLOR_RGB2GRAY)
            piece_cv = cv2.cvtColor(np.array(piece_img), cv2.COLOR_RGB2GRAY)

            bg_img_h, bg_img_w = bg_cv.shape[:2]
            piece_h, piece_w = piece_cv.shape[:2]

            # DOM 宽度与像素宽度比例
            bg_dom_w = bg_elem.size["width"]
            scale_x = bg_img_w / bg_dom_w if bg_dom_w > 0 else 1.0

            # 使用 Canny 边缘增强轮廓
            bg_edges = cv2.Canny(bg_cv, 50, 150)
            piece_edges = cv2.Canny(piece_cv, 50, 150)

            # 模板匹配
            res = cv2.matchTemplate(bg_edges, piece_edges, cv2.TM_CCOEFF_NORMED)

            # 排除当前拼图所在区域，避免匹配到自身
            current_x_px = int(current_x_dom * scale_x)
            exclude_margin = piece_w
            exclude_x_start = max(0, current_x_px - exclude_margin)
            exclude_x_end = min(bg_img_w, current_x_px + piece_w + exclude_margin)
            res[:, exclude_x_start:exclude_x_end] = -1

            # 获取匹配结果
            _, max_val, _, max_loc = cv2.minMaxLoc(res)
            gap_x_px = max_loc[0]

            self.log(
                f"OpenCV 模板匹配: 得分={max_val:.3f}, 当前X(dom)={current_x_dom:.1f}, "
                f"gap_x_px={gap_x_px}"
            )

            # 像素坐标换算到 DOM 坐标
            gap_x_dom = gap_x_px / scale_x

            # 拖动距离 = 缺口位置 - 当前拼图块位置
            distance_dom = gap_x_dom - current_x_dom

            # 略微多拖一点，防止偏左
            distance_dom += 2

            # 限制合理范围
            bg_width_dom = bg_elem.size["width"]
            distance_dom = max(30, min(distance_dom, bg_width_dom - 30))

            self.log(
                f"OpenCV 计算结果: gap_x_dom={gap_x_dom:.1f}, "
                f"distance_dom={distance_dom:.1f}"
            )

            return int(distance_dom)
        except Exception as exc:  # noqa: BLE001
            self.log(f"OpenCV 图像识别异常: {exc}")
            return 0

    def _drag_slider(self, slider, distance: int) -> None:
        """执行滑块拖动操作：单次快速拖动，不分步，避免视觉卡顿。"""
        assert self.driver is not None
        actions = ActionChains(self.driver)
        
        # 点击并按住滑块
        actions.click_and_hold(slider).perform()

        # 单次直接拖到目标位置（避免多步造成卡顿感）
        actions.move_by_offset(distance, 0).perform()
        # 极短停顿，模拟人手松开前的轻微停顿
        time.sleep(0.02)

        # 轻微回拉（很快完成，不显著影响总时长）
        actions.move_by_offset(-1, 0).perform()
        time.sleep(0.01)

        # 释放
        actions.release().perform()
        self.log("已完成滑块拖动操作")

    def fetch_user_points(self) -> float | None:
        """获取用户积分信息。"""
        assert self.driver is not None
        
        try:
            self.log("正在访问积分页面...")
            self.driver.get("https://basic.smartedu.cn/user/myIncentives")
            
            # 等待页面加载
            wait = WebDriverWait(self.driver, 15)
            time.sleep(2)  # 额外等待确保页面完全渲染
            
            # 尝试多种策略定位积分元素
            points_text = None
            strategies = [
                # 策略1: 通过类名定位（常见的积分显示类名）
                (By.CSS_SELECTOR, ".score-value, .points-value, .integral-num, .score-num"),
                # 策略2: 通过包含"我的积分"的父元素查找
                (By.XPATH, "//*[contains(text(), '我的积分')]/following::*[1]"),
                # 策略3: 查找包含数字的大号文本元素（积分通常显示较大）
                (By.XPATH, "//*[contains(@class, 'score') or contains(@class, 'point') or contains(@class, 'integral')]//text()[normalize-space()]"),
                # 策略4: 查找所有可能包含积分的元素
                (By.XPATH, "//*[contains(text(), '.') and string-length(normalize-space(text())) < 10]"),
            ]
            
            for idx, (by, selector) in enumerate(strategies, 1):
                try:
                    self.log(f"尝试策略 {idx}: {selector}")
                    elements = self.driver.find_elements(by, selector)
                    
                    for elem in elements:
                        try:
                            text = elem.text.strip()
                            if text:
                                # 使用正则表达式提取数字（支持整数和小数）
                                match = re.search(r'\d+\.?\d*', text)
                                if match:
                                    points_text = match.group()
                                    self.log(f"策略 {idx} 找到积分文本: {points_text} (原始: {text})")
                                    break
                        except Exception:  # noqa: BLE001
                            continue
                    
                    if points_text:
                        break
                        
                except Exception as exc:  # noqa: BLE001
                    self.log(f"策略 {idx} 失败: {exc}")
                    continue
            
            # 如果所有策略都失败，尝试从页面源码中提取
            if not points_text:
                self.log("所有定位策略失败，尝试从页面源码提取...")
                page_source = self.driver.page_source
                
                # 查找"我的积分"附近的数字
                patterns = [
                    r'我的积分[^0-9]*(\d+\.?\d*)',
                    r'当前积分[^0-9]*(\d+\.?\d*)',
                    r'积分[^0-9]*(\d+\.?\d*)',
                ]
                
                for pattern in patterns:
                    match = re.search(pattern, page_source)
                    if match:
                        points_text = match.group(1)
                        self.log(f"从页面源码提取到积分: {points_text}")
                        break
            
            # 保存调试截图
            try:
                screenshot_path = f"debug_points_{datetime.now().strftime('%Y%m%d_%H%M%S')}.png"
                self.driver.save_screenshot(screenshot_path)
                self.log(f"已保存调试截图: {screenshot_path}")
            except Exception:  # noqa: BLE001
                pass
            
            # 转换为浮点数
            if points_text:
                try:
                    points = float(points_text)
                    self.log(f"成功获取积分: {points}")
                    return points
                except ValueError as exc:
                    self.log(f"积分文本转换失败: {points_text} -> {exc}")
                    return None
            else:
                self.log("未能找到积分信息")
                return None
                
        except Exception as exc:  # noqa: BLE001
            self.log(f"获取积分异常: {exc}")
            return None

    @staticmethod
    def _generate_track(distance: int) -> list[int]:
        """保留接口，返回整体距离（兼容调用，不再分步使用）。"""
        return [distance]

