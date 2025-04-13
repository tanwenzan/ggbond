package main

import "fmt"

func main() {
	// sqlx
	QueryByDepart()
	QueryMaxSalaryEmployee()

	// books
	QueryByBook()

	post, err := GetMostCommentedPost()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(post)
	}

	comments, err := GetUserPostsWithComments(1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(comments)
	}
}
