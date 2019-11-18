package workscheduler
// copy from http://nil.csail.mit.edu/6.824/2018/notes/gopattern.pdf

// version1:
//func Schedule(servers []string, numTask int, call func(srv string, task int)) {
//    idle := make(chan string, len(servers))
//    for _, srv := range servers {
//        idle <- srv
//    }
//
//    for task := 0; task < numTask; task ++ {
//        srv := <- idle
//        go func(t int) {
//            call(srv, t)
//            idle <- srv
//        }(task)
//    }
//    for i:= 0; i < len(servers); i++ {
//        <- idle
//    }
//}

// version2: 有点mapreduce 的味道了
func Schedule(servers chan string, numTask int, call func(srv string, task int) bool) {
    // Think carefully beforeintroducing unbounded queuing
    work := make(chan int, numTask)
    done := make(chan bool)
    exit := make(chan bool)

    runTasks := func(srv string) {
        for task := range work {
            if call(srv, task) {
                done <- true
            } else {
                work <- task
            }
        }
    }

    go func() {
        for {
            select {
            case srv := <- servers:
                go runTasks(srv)
            case <- exit:
                return
            }
        }
    }()

    for task := 0; task < numTask; task++ {
        work <- task
    }
    // Close a channel to signalthat no more values will be sent
    //
    for i := 0; i < numTask; i++ {
        <- done
    }
    close(work)
    exit <- true

}
