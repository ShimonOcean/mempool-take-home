# Shimon Islam K--- L--s Take Home Assessment

(Anonymized Company Name to prevent plagiarism)

Done in Golang, ran in the /src directory with

```go run main.go``` 

If priority is calculated as FeePerGas * Gas for each transaction, then I thought the best way to implement the system described would be through a priority queue. Each transaction is stored in this priority queue and acts like a min heap, where the priority for each transaction are determined by FeePerGas * Gas.

After storing each transaction (up to 5000) and popping transactions off the queue with the least priority once reaching the limit, I print the 5000 stored transactions into the prioritized-transactions.txt file. 

I have added a couple of tests in main_test.go, and admittedly did not thoroughly check every possible edge case but I believe I have covered some of the important ones. 

Please feel free to send me any feedback for improvements or different ways I could have solved this problem! shimon.islam88@gmail.com
