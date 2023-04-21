#!/bin/python
import colorama

# Initialize colorama
colorama.init()

# Define a dictionary that maps keywords to colors
keyword_colors = {
    "INFO": colorama.Fore.YELLOW,
    "ERROR": colorama.Fore.RED,
    "WARNING": colorama.Fore.MAGENTA,
    "CODE": colorama.Style.DIM 
}

# Define a function to print messages with colored keywords
def print_message(keyword, message):

    """Prints messages with colored keywords."""

    color = keyword_colors.get(keyword, "")  # Get the color for the keyword from the dictionary
    colored_keyword = f"{color}{keyword}:{colorama.Style.RESET_ALL}" if color else keyword  # Add color to the keyword if it exists in the dictionary
    print(f"{colored_keyword} {message}")

