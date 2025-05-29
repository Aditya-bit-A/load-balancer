# Load Balancer

A simple, lightweight Layer 7 load balancer implemented in Go. It distributes incoming HTTP requests to multiple backend servers using Round Robin or Least Connections strategy.

## Features

- Supports Round Robin and Hashing algorithms
- Uses consistent hashing interally.
- Graceful handling of backend server failures
- Simple and clean Go implementation
- Logs incoming requests and backend assignments

## Getting Started

### Prerequisites

- Go 1.18 or higher

### Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/Aditya-bit-A/load-balancer.git
   cd load-balancer
