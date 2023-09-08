# R.A.L.P.H. - Rapid Automated Load Performance Harness

## Description
RALPH is a load testing tool written in Go. It's designed to test the performance of websites using various load testing strategies. It can simulate simple tests, stress tests, spike tests, endurance tests, and ramp-up tests. Upon completion, RALPH saves the results as a CSV file and uploads it to Google Cloud Storage.

Certainly! Installation and build instructions are important aspects of any README, especially for projects where users might be cloning and building the code themselves. I'll provide a section for installation and build instructions that you can include in your README:

---

## Installation and Build Instructions

### Prerequisites:

1. **Go (Golang)**: Ensure you have Go installed. If not, you can download and install it from [the official Go website](https://golang.org/dl/).

2. **Google Cloud SDK**: If you're planning to use the Google Cloud Storage upload functionality, make sure you have the Google Cloud SDK installed. You can get it from [the official Google Cloud SDK website](https://cloud.google.com/sdk/docs/install).

3. **Environment Variable**: Ensure the `GOOGLE_APPLICATION_CREDENTIALS` environment variable is set up if you're using Google Cloud Storage. This should point to the path of your service account key. [Learn more here](https://cloud.google.com/docs/authentication/getting-started).

### Installation:

1. **Clone the repository**:

```bash
git clone https://github.com/Chaos-Camp/RAPLH.git
cd RAPLH
```

2. **Fetch dependencies**:

Your project seems to be using external packages, so ensure you fetch all required dependencies:

```bash
go mod tidy
```

### Build:

1. **Compile the source code**:

This will generate an executable named `loadtester` (or `loadtester.exe` on Windows).

```bash
go build -o loadtester main.go
```

### Running:

Once you've built the application, you can run it directly:

```bash
./loadtester -url=https://www.test-site.com -type=simple -iterations=3
```

(Note: On Windows, the executable will be `loadtester.exe`)

---

This section should give users clear guidance on how to get your project up and running on their local machines. Adjust the placeholders (`your-username`, `your-repo-name`, etc.) to match your actual repository details.

## Usage

Set up the environment variable for GCS authentication:

```bash
export GOOGLE_APPLICATION_CREDENTIALS=<path_to_service_account_json>
```

### Options:

- `-url`: The URL you want to test. Default is "https://www.example.com".
- `-type`: Type of load test (e.g., simple, stress, spike, endurance, ramp_up). Default is "ramp_up".
- `-iterations`: Number of iterations for a simple test. Default is 1.
- `-concurrentRequests`: Number of concurrent requests for a stress test. Default is 5.
- `-spikes`: Number of spikes for a spike test. Default is 3.
- `-spikeInterval`: Interval (in seconds) between spikes. Default is 1.0 second.
- `-duration`: Duration (in seconds) for the endurance test. Default is 300 seconds (5 minutes).
- `-maxUsers`: Maximum number of users for the ramp-up test. Default is 10.
- `-rampUpPeriod`: Ramp-up period (in seconds). Default is 10 seconds.
- `-bucket`: Google Cloud Storage bucket name where the results will be uploaded. Default is "my-bucket".

### Examples:

1. **Simple Test** on "https://www.test-site.com" with 3 iterations:

```bash
go run main.go -url=https://www.test-site.com -type=simple -iterations=3
```

2. **Stress Test** on "https://www.stress-site.com" with 10 concurrent requests:

```bash
go run main.go -url=https://www.stress-site.com -type=stress -concurrentRequests=10
```

3. **Spike Test** on "https://www.spike-site.com" with 5 spikes and an interval of 2 seconds between them:

```bash
go run main.go -url=https://www.spike-site.com -type=spike -spikes=5 -spikeInterval=2
```

4. **Endurance Test** on "https://www.endurance-site.com" lasting 600 seconds (10 minutes):

```bash
go run main.go -url=https://www.endurance-site.com -type=endurance -duration=600
```

5. **Ramp Up Test** on "https://www.rampup-site.com" with a maximum of 15 users and a total ramp-up period of 15 seconds:

```bash
go run main.go -url=https://www.rampup-site.com -type=ramp_up -maxUsers=15 -rampUpPeriod=15
```

6. **Save Results** to a specific GCP bucket after running any of the above tests:

```bash
go run main.go -url=https://www.example-site.com -type=spike -bucket=my-special-bucket
```

---

These examples provide a comprehensive overview of how to use the tool for various types of tests. Adjusting the provided URLs or other parameters can further tailor the tests to specific requirements.

After the test, results will be saved in a CSV file named after the sanitized URL followed by `_results.csv`. This CSV file will then be uploaded to the specified Google Cloud Storage bucket.

## Contributing

Please ensure to follow Go's coding conventions and best practices when submitting a pull request.

## License

This project is open-sourced and available for all. Do make sure to provide proper attribution if using or referencing the code.


