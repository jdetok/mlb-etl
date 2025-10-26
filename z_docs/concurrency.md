# outline concrrent needs for box score etl
- HIGH LEVEL:
    - set a range of years (seasons)
    - query the database to receive the full list of game ids from those seasons
        - structured in a map with a game id mapped to a season
        - NEEDS TO OCCUR FIRST (boxscore fetch relies on game id)
    - split the slice of gameid:season maps into many chunks of slices of slices
    - <b>CHANNELS NEED TO BE CREATED HERE</b>
        - need channels for errors for the goroutines to write to and then the logger picks up
            - dberrs chan
            - httperrs chan
            - success chan
    - GO loop through the chunks (goroutine created)
        - GO loop through each gameid:season map in the chunk
            - call MakeMultiTableETL to make a multi table ETL object (ETL)
            - call ETL.ExtractData to run the http request & convert to the appropriate struct (must implement the ETLProcess interface)
            - call the interface methods on the data
            - loop through each PGTarget (each target table)
                - call BuildAndInsert to build the insert statement and execute it in the database
        - when a season:map finishes, write to the success chan
    - GO func for reading the chans and logging them


