package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// task1 - 定义与数据库字段映射的结构体

type Employee struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Salary     int    `db:"salary"`
}

func QueryByDepart() {
	// 假设已建立数据库连接
	db := sqlx.MustConnect("mysql", "user:password@tcp(localhost:3306)/dbname")

	// 查询技术部员工
	var techEmployees []Employee
	err := db.Select(&techEmployees,
		`SELECT id, name, department, salary 
         FROM employees 
         WHERE department = "技术部"`)

	if err != nil {
		panic(err)
	}

	fmt.Printf("技术部共 %d 名员工:\n", len(techEmployees))
	for _, emp := range techEmployees {
		fmt.Printf("ID:%d 姓名:%s 薪资:%d\n", emp.ID, emp.Name, emp.Salary)
	}
}

// task1 - 最高薪资员工查询实现

func QueryMaxSalaryEmployee() {

	db := sqlx.MustConnect("mysql", "user:password@tcp(localhost:3306)/dbname")

	var topEmployee Employee

	// 查询薪资最高的员工[9,11](@ref)
	err := db.Get(&topEmployee,
		`SELECT id, name, department, salary 
         FROM employees 
         ORDER BY salary DESC 
         LIMIT 1`)

	if err != nil {
		panic(err)
	}

	fmt.Printf("\n最高薪资员工:\nID:%d 姓名:%s 部门:%s 薪资:%d\n",
		topEmployee.ID, topEmployee.Name,
		topEmployee.Department, topEmployee.Salary)
}

// task2 - 实现类型安全映射

type Book struct {
	ID     int    `db:"id"`
	Title  string `db:"title"`
	Author string `db:"author"`
	Price  int    `db:"price"`
}

func QueryByBook() {
	db := sqlx.MustConnect("mysql", "user:password@tcp(localhost:3306)/dbname")
	var books []Book

	err := db.Get(&books, "select id,title,author,price from books where price > 50")

	if err != nil {
		panic(err)
	}

	for _, v := range books {
		fmt.Printf("\n书籍信息:\nID:%d 书名:%s 作者:%s 价格:%d\n",
			v.ID, v.Title, v.Author, v.Price)
	}
}
