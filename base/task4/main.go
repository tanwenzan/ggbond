package main

import "fmt"

func main() {
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
