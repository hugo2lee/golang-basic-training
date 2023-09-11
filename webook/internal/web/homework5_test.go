/*
 * @Author: hugo2lee hugo2lee@gmail.com
 * @Date: 2023-09-10 23:58
 * @LastEditors: hugo2lee hugo2lee@gmail.com
 * @LastEditTime: 2023-09-11 18:22
 * @FilePath: /geektime-basic-go/webook/internal/web/homework5_test.go
 * @Description:
 *
 * Copyright (c) 2023 by hugo, All Rights Reserved.
 */
package web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	svcmocks "gitee.com/geekbang/basic-go/webook/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_LoginJWT(t *testing.T) {
	const loginURL = "/users/login"

	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request
		// 预期响应
		wantCode int
		wantBody string
		wantUID  int64
		jwtParse func(t *testing.T, tokenStr string) int64
	}{
		// TODO: Add test cases.
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().
					Login(gomock.Any(), "123@qq.com", "hello@world123").
					Return(domain.User{
						Id: 3,
					}, nil)
				codesvc := svcmocks.NewMockCodeService(ctrl)
				return usersvc, codesvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email":"123@qq.com","password":"hello@world123"}`))
				req, err := http.NewRequest(http.MethodPost, loginURL, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "登录成功",
			wantUID:  3,
			jwtParse: func(t *testing.T, tokenStr string) int64 {
				uc := UserClaims{}
				token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (any, error) {
					return JWTKey, nil
				})
				assert.NoError(t, err)
				assert.True(t, token.Valid)
				return uc.Id
			},
		},
		{
			name: "请求参数错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				usersvc := svcmocks.NewMockUserService(ctrl)
				codesvc := svcmocks.NewMockCodeService(ctrl)
				return usersvc, codesvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email":"123@qq.com",,,,"password":"hello@world123"}`))
				req, err := http.NewRequest(http.MethodPost, loginURL, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusBadRequest,
			wantBody: "",
			wantUID:  0,
			jwtParse: func(t *testing.T, tokenStr string) int64 {
				return 0
			},
		},
		{
			name: "用户名或者密码不正确",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().
					Login(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(domain.User{}, service.ErrInvalidUserOrPassword)
				codesvc := svcmocks.NewMockCodeService(ctrl)
				return usersvc, codesvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email":"123@qq.com","password":"hello@world123"}`))
				req, err := http.NewRequest(http.MethodPost, loginURL, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "用户名或者密码不正确，请重试",
			wantUID:  0,
			jwtParse: func(t *testing.T, tokenStr string) int64 {
				return 0
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// 准备 mock
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			usersvc, codesvc := tt.mock(ctrl)
			han := NewUserHandler(usersvc, codesvc)

			// 准备 请求和响应
			server := gin.Default()
			han.RegisterRoutes(server)
			req := tt.reqBuilder(t)

			// 执行
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			assert.Equal(t, tt.wantCode, recorder.Code)
			assert.Equal(t, tt.wantBody, recorder.Body.String())

			// 验证
			tokenStr := recorder.Result().Header.Get("x-jwt-token")
			assert.Equal(t, tt.wantUID, tt.jwtParse(t, tokenStr))
		})
	}
}
