package main

import (
	"context"
	math "happyladysauce/kitex_gen/math"
)

// MathServiceImpl implements the last service interface defined in the IDL.
type MathServiceImpl struct{}

// Add implements the MathServiceImpl interface.
func (s *MathServiceImpl) Add(ctx context.Context, req *math.AddRequest) (resp *math.AddResponse, err error) {
	// TODO: Your code here...
	resp = &math.AddResponse{
		Result: req.Left + req.Right,
	}
	return
}

// Sub implements the MathServiceImpl interface.
func (s *MathServiceImpl) Sub(ctx context.Context, req *math.SubRequest) (resp *math.SubResponse, err error) {
	// TODO: Your code here...
	resp = &math.SubResponse{
		Result: req.Left - req.Right,
	}
	return
}
