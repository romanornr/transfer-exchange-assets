// Copyright (c) 2021 Romano
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package queue

var QueueDispatcher *Dispatcher

func InitQueueDispatcher() {
	QueueDispatcher = NewDispatcher(4)
	QueueDispatcher.Run()
}

func Push(job Queuable) {
	JobQueue <- job
}
