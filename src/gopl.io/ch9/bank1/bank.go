// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 261.
//!+

// Package bank provides a concurrency-safe bank with one account.
package bank

import "fmt"

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance: // select 多路复用时, 解决 Channel 发送消息对方没有准备好而堵塞的问题
		// 这里的 balances <- balance 只有在 balances 准备好接收时才会调用, 也就是说调用 Balance() 函数才会执行到此case语句
		// 所以不会出现0
			fmt.Println("balance: ", balance)
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}

//!-
