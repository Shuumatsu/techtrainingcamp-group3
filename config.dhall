let log_dir = "./temp" in 
    {
        zap_log_file = log_dir ++ "/stdout.log",
        gin_log_file = log_dir ++ "/gin.log"
    }
