package main 

import ("fmt"; "context"; "errors"; "net/http"; "time"; "os"; "os/signal"; "syscall")

func (app *application) serve() error{
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s:=<-quit 

		app.logger.Printf("Shutting down server.. %s\n", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
		  
	}()
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed){
		return err 
	}

	err = <-shutdownError
	if err != nil{
		return err 
	}

	app.logger.Printf("Stopped server..")
	return nil
}
