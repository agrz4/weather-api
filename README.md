This project helps you build a fast and reliable Weather API using Go. It uses a popular web framework called chi for handling requests and Redis to quickly store and retrieve weather data, making it very efficient.

Key things it does:

- Gets Live Weather: Fetches current weather for any city you ask for.
- Super Fast with Caching: Saves weather data in Redis so it doesn't have to ask the weather service every time, making it much quicker and cheaper.
- Easy to Use: Has simple web addresses like /health (to check if it's working) and /v1/weather (to get weather data).
- Smart City Handling: Understands different ways you might type a city name (like "New York" or "New-York").
- Flexible: You can easily change settings like API keys and server details.
