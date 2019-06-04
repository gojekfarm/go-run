This package allows you to set-up a background processing system for your golang application.

It consists of three main parts
* Dispatcher 
    - Holds the worker pool and job queue.
    - Acts as an interface to queue jobs.
* Worker
    - These are goroutines that are constantly running and polling their internal 
      job queue for jobs and they execute them as soon as possible, if a job errors
      out or panics, they retry the job depending on the configuration set.
    - The number of workers running at the time can be set through the configuration
* Job
    - Acts as an interface which is used in the workers and dispatcher.
    - To make something act as a job, all one needs to do is implement the method
      Execute() which returns an error.
  