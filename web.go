package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", authMiddleware(handler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	var currentUser User

	// 컨텍스트에서 값을 가져옴
	if v := r.Context().Value("current_user"); v == nil {
		// "current_user"가 존재하지 않으면 401 에러 리턴
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	} else {
		u, ok := v.(User)
		if !ok {
			// 타입이 User가 아니면 401 에러 리턴
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}

		currentUser = u
	}

	fmt.Fprintf(w, "Hi I am %s", currentUser.Name)
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 현재 세션 정보를 확인하여 currentUser 생성
		currentUser, err := getCurrentUser(r)
		if err != nil {
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}

		// 기본 컨텍스트(r.Context())에 current_user 값을 담은 새로운 컨텍스트 생성
		ctx := context.WithValue(r.Context(), "current_user", currentUser)

		// 기존 http.Request의 컨텍스트를 변경
		nextRequest := r.WithContext(ctx)

		// 다음 handlerFunc 호출
		next(w, nextRequest)
	}
}

func getCurrentUser(r *http.Request) (User, error) {
	return User{Name: "Jaehue"}, nil
}

type User struct {
	Name string
}
