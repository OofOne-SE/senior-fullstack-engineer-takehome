import fs from 'fs';

const logger = console

const config = {
    lineDelimiter: "\n",
    columnDelimiter: "\t",
    weatherDataFilePath: "./resources/weather.dat",
    endpoint: "http://localhost:8080/api/v1/weather",
    debounceTimeInMilliseconds: 2000
}


const loadWeatherData = (lineDelimiter, columnDelimiter, filePath) => {
    console.log(filePath);
    const weatherDataFile = fs.readFileSync(filePath, 'utf8')
    return weatherDataFile.split(lineDelimiter).filter(line => !!line).map(line => {
        const [date, temperature, humidity] = line.split(columnDelimiter)
        return {
            date: date, humidity: parseFloat(humidity), temperature: parseFloat(temperature)
        }
    })
}

const isWeatherDataValid = (weatherData) => {
    if (!weatherData.length) {
        logger.warn('Weather data is empty.')
        return false
    }
    const isRecordInvalid = record => {
        const isInvalid = !record.date || isNaN(record.humidity) || isNaN(record.temperature)
        if (isInvalid) {
            logger.warn(`Invalid record: ${JSON.stringify(record)}`)
        }
        return isInvalid
    }

    if (weatherData.some(isRecordInvalid)) {
        logger.warn('Weather data contains invalid records.')
        return false
    }

    return true
}

const getWeatherDataAndSendToApi = async ({ lineDelimiter, columnDelimiter, weatherDataFilePath }) => {
    const weatherData = loadWeatherData(lineDelimiter, columnDelimiter, weatherDataFilePath)
    if (!isWeatherDataValid(weatherData)) {
        logger.error('Weather data is invalid, aborting...')
        return
    }
    for (const item of weatherData) {
        await new Promise(resolve => setTimeout(resolve, config.debounceTimeInMilliseconds))
        logger.debug(`Sending data: ${item.date}`)
        fetch(config.endpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(item)
        })
            .then(response => response.json())
            .catch(error => logger.error(error))
    }
}

getWeatherDataAndSendToApi(config);
