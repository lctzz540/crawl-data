# Grade 10 Entrance Exam Data Crawler and Graph Plotter

This repository contains a data crawler and graph plotter specifically designed for the 10th grade entrance exam data. The crawler extracts data from the exam results and generates graphical representations of the scores for each subject. This tool aims to provide insights and visualizations for better understanding and analysis of the exam data.

## Features

- Data crawling: The crawler fetches the exam data from a specified source, such as a website or CSV file.
- Subject-wise graph plotting: The tool generates individual bar charts for each subject, visualizing the distribution of scores.
- Customizable chart settings: Various chart options, such as titles, subtitles, axis labels, and bar colors, can be configured to suit specific needs.
- Automatic data processing: The tool automatically processes the exam data and calculates the score distribution for each subject.

## Installation

To use the Grade 10 Entrance Exam Data Crawler and Graph Plotter, follow these steps:

1. Clone the repository to your local machine:

   ```bash
   git clone https://github.com/lctzz540/crawldata.git
   ```

2. Install package

   ```bash
   go get
   ```

3. Usage

   To update data, let run:

   ```bash
   go run . -update
   ```

   To plotting data, let run:

   ```bash
   go run . -plot
   ```
