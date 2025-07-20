# AI-Powered Debugging System
# TeaThis - Advanced Code Analysis Platform

TeaThis is a comprehensive code analysis platform that helps developers diagnose issues in their code through static analysis, error log parsing, and knowledge graph-based error localization.

![TeaThis Dashboard](docs/dashboard.png)

## Features

- **Repository Management**: Clone, index, and manage Git repositories
- **Static Code Analysis**: Analyze code for potential issues and vulnerabilities
- **Error Log Parsing**: Extract and interpret error messages from build and runtime logs
- **Code Knowledge Graph**: Visualize relationships between components in your codebase
- **Root Cause Diagnosis**: Identify the underlying causes of errors using graph analysis
- **Modern UI**: Beautiful, responsive user interface built with React and Material UI

## Tech Stack

### Backend
- Go 1.24
- Gin Web Framework
- go-git for Git operations
- Neo4j (optional) for advanced graph storage

### Frontend
- React 18
- Material UI
- D3.js for graph visualization
- React Router for navigation

### DevOps
- Docker & Docker Compose
- GitHub Actions for CI/CD
- Makefile for common tasks

## Getting Started

### Prerequisites
- Go 1.24 or later
- Node.js 16 or later
- Docker and Docker Compose (optional)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/teathis/codeanalyzer.git
   cd codeanalyzer
   ```

2. Install dependencies:
   ```bash
   make install-deps
   ```

3. Build the application:
   ```bash
   make build
   ```

4. Run the application:
   ```bash
   make run
   ```

The application will be available at http://localhost:8080

### Development Setup

1. Start the backend in development mode:
   ```bash
   make backend-dev
   ```

2. Start the frontend development server:
   ```bash
   make frontend-dev
   ```

The frontend will be available at http://localhost:3000 with hot reloading enabled.

### Docker Setup

1. Build and run with Docker:
   ```bash
   make docker-run
   ```

2. Or use Docker Compose for a complete development environment:
   ```bash
   make dev-up
   ```

## Project Structure
This repository contains an AI-powered debugging system that uses computer vision, reinforcement learning, and knowledge graphs to diagnose and resolve cursor-related issues.

## Features

- Real-time cursor tracking with depth-aware segmentation
- Knowledge graph-based diagnostic reasoning
- Reinforcement learning agent for action planning
- Web-based monitoring and control interface

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Neo4j database server
- Testing packages (for development)

### Installation

1. Clone the repository:

#### AI Debugging Agent Workflow

1. **Clone the repo**
   ```bash
   git clone https://github.com/apache/spark.git
   ```

2. **Run the AI debugging agent**
   ```bash
   python /root/GolandProjects/AI/auto_debug_workflow.py /path/to/spark
   ```

3. **Review the report**
   - Check `debug_report.txt` for summary of fixes and recommendations.

4. **Deploy fixes**
   - After confirmation, push changes to production.
