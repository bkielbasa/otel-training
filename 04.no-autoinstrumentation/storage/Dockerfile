# Use an official lightweight version of the Python image.
FROM python:3.9-slim

# Set the working directory in the Docker container.
WORKDIR /app

# Copy the local code to the container's working directory.
COPY . /app

# Install the Python dependencies.
RUN pip install --no-cache-dir Flask requests boto3

# Make port 5000 available to the world outside this container
# as Flask runs on port 5000 by default.
EXPOSE 5000

# Define environment variable for Flask to run in production mode.
ENV FLASK_ENV=production
ENV FLASK_APP=main.py

# Use the flask run command to start your application.
CMD ["flask", "run", "--host=0.0.0.0"]

