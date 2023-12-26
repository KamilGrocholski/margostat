# How does it work

## Every n minutes (SCRAPER_INTERVAL_IN_MINUTES set inside the .env) the program runs:
1. Fetch online players from every world
2. Insert the fetched data into the database
3. Select and parse the data from the database
4. Create or replace html files for each world and index inside the static dir 

## And the server serves the created html files

# How to run

## Set .env using the .env.example file

## Start docker
```bash
docker-compose up --build
```

## To see the result, go to:
```bash
http://localhost:8080
```

<p>When the database is empty, there will be no html files to serve, so wait until the first job is done</p>
