import curses
import time
import random
import sys
import argparse


def matrix_effect(stdscr, lines, max_active_lines=1000, freedom=True, speed=1, ranNum=100):
    curses.curs_set(0)  # Hide the cursor
    stdscr.nodelay(1)   # Make getch() non-blocking
    stdscr.timeout(100) # Refresh every 100ms

    # Enable colors
    curses.start_color()
    color_pairs = [
        (curses.COLOR_GREEN, curses.COLOR_BLACK),
        (curses.COLOR_BLUE, curses.COLOR_BLACK),
        (curses.COLOR_RED, curses.COLOR_BLACK),
        (curses.COLOR_YELLOW, curses.COLOR_BLACK),
        (curses.COLOR_MAGENTA, curses.COLOR_BLACK)
    ]

    # Initialize color pairs
    for i, (fg, bg) in enumerate(color_pairs, start=1):
        curses.init_pair(i, fg, bg)

    # Get screen size
    height, width = stdscr.getmaxyx()

    # Initialize positions for all lines in the file
    positions = [
        {
            "row": 0 if freedom else height,  # Start at the top if 'freedom' is enabled, else at the bottom
            "col": 0 if not freedom else random.randint(0, max(0, width - len(line))),  # Random or forced alignment
            "speed": speed,  # Use the user-defined speed
            "text": f"Line: {i + 1} | Text: {line}",
            "color": ((i // ranNum) % len(color_pairs)) + 1  # Assign color based on index
        }
        for i, line in enumerate(lines)
    ]

    visible_lines = []  # Rolling window of active lines
    index = 0  # Start with the first line

    while True:
        stdscr.clear()

        # Ensure dimensions are recalculated in case of terminal resize
        height, width = stdscr.getmaxyx()

        # Add lines to the rolling window up to max_active_lines
        if len(visible_lines) < max_active_lines:
            visible_lines.append(positions[index])
            index = (index + 1) % len(positions)

        for pos in visible_lines:
            # Truncate text to fit the screen width
            truncated_text = pos["text"][:width]

            # Use the pre-assigned color for this line
            color_index = pos["color"]

            # Display the line only if within screen bounds
            if 0 <= pos["row"] < height:
                try:
                    stdscr.addstr(int(pos["row"]), pos["col"], truncated_text, curses.color_pair(color_index))
                except curses.error:
                    pass

            # Move the line position depending on the 'freedom' argument
            if freedom:
                pos["row"] += pos["speed"]  # Move downward if falling from the top
            else:
                pos["row"] -= pos["speed"]  # Move upward if starting from the bottom

        # Remove lines that go off-screen
        if freedom:
            visible_lines = [pos for pos in visible_lines if pos["row"] < height]
        else:
            visible_lines = [pos for pos in visible_lines if pos["row"] >= 0]

        # Refresh the screen
        stdscr.refresh()

        # Break if the user presses a key
        if stdscr.getch() != -1:
            break

        # Slow down the loop for effect
        time.sleep(0.1)

def main(file_name, max_active_lines, freedom, speed, ranNum):
    with open(file_name, 'r') as f:
        lines = [line.strip() for line in f.readlines()]

    curses.wrapper(lambda stdscr: matrix_effect(stdscr, lines, max_active_lines, freedom, speed, ranNum))


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Matrix-style scrolling text effect.")
    parser.add_argument("file_name", type=str, help="Path to the file containing text.")
    parser.add_argument("--max_active_lines", type=int, default=1000, help="Maximum number of lines visible at a time.")
    parser.add_argument("--freedom", action="store_true", default=False, help="Make text fall from the top.")
    parser.add_argument("--speed", type=int, default=1, help="Set the scroll speed (higher is faster).")
    parser.add_argument("--group", type=int, default=100, help="Set the grouping number.")
    args = parser.parse_args()

    main(args.file_name, args.max_active_lines, args.freedom, args.speed, args.group)
