# Auto Meeting Scheduler
The purpose of the program is to automatically suggest meeting times between users who desire to meet, taking into account already existing meetings & the availability of each user. 

Availability of users is specified in a JSON file: `samples/source.json` and it has the following structure:
```
[
     {
        "mettings": [
            [
                "9:00", //start of the meeting
                "10:30" //end of the meeting
            ],
            [
                "12:00",
                "13:00"
            ],
            [
                "16:00",
                "18:00"
            ]
        ],
        "dailyBounds": [
            "8:00", //start of work day
            "18:00" //end of work day
        ]
    }
]
```
## Running the program
To run the program use `go run main.go [desired duration of the meeting in minutes]`. For example, the following code will suggest 15 minute slots of possible meeting times between users: 
```
go run main.go 15
```

## License
The license is specified in the file [LICENSE](https://github.com/JoelD7/auto-scheduler/blob/master/LICENSE).