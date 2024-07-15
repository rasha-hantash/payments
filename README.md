# chariot-takehome
# chariot-takehome


Setup 
Run brew install go-task
Run task setup
Run task proto:gen:api
Run docker compose up -d  --no-deps --build
Run task migrate:up:local


Cleanup
Run docker compose down 


# TODO for today
- complete readme
    - write out picture of architecture, justification for logging, 
- make sure context is passed from gateway to 
- (yes) add appropriate history trigger and custom_id trigger 
- (yes) add logging 
- make sure tests pass and write go tests 
- add in bench mark tests
- extra: add middleware to gateway to check for authentication
-  extra: add rate limiter
- extra: add health monitor 
- extra: add a pr branch for how to do an ach transaction with a 3p like moov 
- run go fmt 
- https://dev.to/stripe/designing-apis-for-humans-object-ids-3o5a
- https://medium.com/@RobertKhou/double-entry-accounting-in-a-relational-database-2b7838a5d7f8 
- https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html
- https://medium.com/@connorbutch/the-importance-of-structured-logging-in-aws-and-anywhere-else-52a4534c53aa
- https://clerk.com/blog/generating-sortable-stripe-like-ids-with-segment-ksuids
- https://segment.com/blog/a-brief-history-of-the-uuid/
- other: https://github.com/paralleldrive/cuid2
