const express = require('express');
const axios = require('axios');
const app = express();
const { Pool } = require('pg');
const PORT = process.env.PORT || 3333;
const fetch = require("node-fetch");

const pool = new Pool({
  user: 'postgres',
  host: process.env.POSTGRES_HOST || 'localhost',
  database: 'temperature_db',
  password: 'your_password',
  port: 5434,
});

const tempHost = process.env.TEMPERATURE_HOST || 'localhost';
const storageHost = process.env.STORAGE_HOST || 'localhost';

app.use(express.json());

app.post('/import', async (req, res) => {
    if (!req.body.addresses) {
        return res.status(400).send({ error: 'Addresses array is missing in the request body.' });
    }

    const addresses = req.body.addresses;
    const weatherDataPromises = addresses.map(async (address) => {
        const url = `http://api.weatherapi.com/v1/current.json?key=5a2d6a9bcdd54cdd97a153506242003&q=${address}`
        try {
            const response = await fetch(url);
            if (!response.ok) { // Check if the request was successful
                throw new Error('Network response was not ok');
            }
            const data = await response.json();

            const windSpeed = data.current.wind_kph;
            const windDirection = data.current.wind_dir;
            const location = data.location.name;
            const localtime = data.location.localtime;

            const insertQuery = `
                INSERT INTO wind_infos (wind_speed, wind_direction, location, "localtime") 
                VALUES ($1, $2, $3, $4)`;

            pool.query(insertQuery, [windSpeed, windDirection, location, localtime], (err, res) => {
              if (err) {
                console.error('Error executing query', err.stack);
              } else {
                console.log('Insert operation successful', res);
              }
              // When done with the connection, release it.
            });
            return data.current;
        } catch (error) {
            console.error('Failed to fetch wind info:', error);
        }

        try {
            const response = await fetch('http://${tempHost}:8080/temperature/${address}')
            if (!response.ok) { // Check if the request was successful
                throw new Error('Network response was not ok');
            }
        } catch (error) {
            console.error('Failed to fetch temperature:', error);
        }
        try {
            const response = await fetch('http://${storageHost}:5000/address/${address}')
            if (!response.ok) { // Check if the request was successful
                throw new Error('Network response was not ok');
            }
        } catch (error) {
            console.error('Failed to fetch temperature:', error);
        }
    });
    
    try {
        const results = await Promise.all(weatherDataPromises);
        res.send(results);
    } catch (error) {
        res.status(500).send({ error: error.message });
    }
});

// Placeholder for /second endpoint
app.get('/second', (req, res) => {
    res.send('Second endpoint response');
});

// Placeholder for /third endpoint
app.get('/third', (req, res) => {
    res.send('Third endpoint response');
});

app.listen(PORT, () => {
    console.log(`Server is running on port ${PORT}`);
});

