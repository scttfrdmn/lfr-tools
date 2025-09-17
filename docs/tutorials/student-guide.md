# Student Guide: Using Your Cloud Computer

This guide helps students connect to and use their cloud computer for class.

## What You Need

Before starting, make sure you have:
- Your **access token** from your teacher
- Your **student ID number**
- A computer with internet connection

## Step-by-Step Setup

### Step 1: Install LFR Tools

**If you have a Mac:**
1. Open the **Terminal** app (find it in Applications → Utilities)
2. Type this command and press Enter:
   ```bash
   brew install lfr
   ```
3. Wait for it to finish installing

**If you have Windows:**
1. Ask your teacher for the Windows installer
2. Download and run the installer
3. Open **Command Prompt** (search for "cmd" in Start menu)

**If you have Linux:**
1. Open your terminal
2. Copy and paste this command:
   ```bash
   curl -L https://github.com/scttfrdmn/lfr-tools/releases/latest/download/lfr_Linux_x86_64.tar.gz | tar xz
   sudo mv lfr /usr/local/bin/
   ```

### Step 2: Set Up Your Access

Your teacher will give you an access token that looks like: `cs101-alice-AbCdEf123`

1. Open your terminal or command prompt
2. Type this command (replace with your actual token and ID):
   ```bash
   lfr connect activate cs101-alice-AbCdEf123 12345
   ```
3. You should see: "✅ Token activated!"

**Important**: You only do this once. The access code is now saved on your computer.

### Step 3: Connect to Your Cloud Computer

Whenever you want to use your cloud computer:

```bash
lfr connect alice
```

Replace "alice" with your actual username.

**What you might see:**
- "Connecting to alice's instance..." ✅ Good!
- "Instance stopped, requesting start..." ⏳ Wait a minute
- "✅ Connected!" 🎉 You're in!

## Using Your Cloud Computer

### Basic Commands

Once connected, you can use these commands:

**See what's in the current folder:**
```bash
ls
```

**Go to your home folder:**
```bash
cd ~
```

**Create a new folder:**
```bash
mkdir my-projects
```

**Create a new file:**
```bash
nano hello.py
```

**Run a Python program:**
```bash
python3 hello.py
```

**Exit back to your own computer:**
```bash
exit
```

### Working with Files

**To edit a file:**
```bash
# Use nano (easier for beginners):
nano myfile.py

# Or use vim (more advanced):
vim myfile.py
```

**To copy files:**
```bash
cp oldfile.py newfile.py
```

**To delete files:**
```bash
rm filename.py
```

**To see what's in a file:**
```bash
cat filename.py
```

### Shared Class Files

Your teacher might set up shared folders that everyone can access:

**See shared files:**
```bash
ls /mnt/efs/shared/
```

**Copy shared files to your folder:**
```bash
cp /mnt/efs/shared/assignment1.py ~/
```

**Submit your work:**
```bash
cp myassignment.py /mnt/efs/submissions/
```

## Common Software

Your cloud computer might have these programs installed:

### Python Programming
```bash
# Run Python:
python3

# Install Python packages:
pip3 install numpy pandas matplotlib

# Run Jupyter Notebook:
jupyter notebook --ip=0.0.0.0 --port=8888 --no-browser
```

### Web Development
```bash
# Create a new website:
npx create-react-app my-website

# Start a web server:
npm start
```

### Data Science
```bash
# Start R:
R

# Use RStudio (if installed):
# Open web browser to: http://your-computer-ip:8787
```

## Important Tips

### Save Your Work Often
- Your cloud computer might turn off automatically to save money
- Always save your files before taking a break
- Use `Ctrl+S` in editors to save

### Internet Connection
- Your cloud computer has fast internet
- Download large files on the cloud computer, not your laptop
- Use `wget` or `curl` to download files

### Getting Unstuck

**If a command doesn't work:**
1. Check for typos in your command
2. Make sure you're in the right folder (`pwd` shows where you are)
3. Try the command again
4. Ask your teacher or TA

**If you can't connect:**
1. Check your internet connection
2. Make sure you typed your username correctly
3. Your computer might be starting up (wait 1-2 minutes)
4. Ask your teacher to check if your computer is running

**If you lost your work:**
1. Check if you saved the file (`ls` to see files)
2. Look in different folders (`cd ~` then `ls`)
3. Check the shared folder (`ls /mnt/efs/shared/`)
4. Ask your teacher - they can help recover files

### Best Practices

**Do:**
- Save your work frequently
- Use descriptive file names like `assignment1_alice.py`
- Put your work in organized folders
- Ask questions when stuck

**Don't:**
- Share your access token with others
- Delete files you didn't create
- Run commands you don't understand
- Panic if something goes wrong (it's usually fixable!)

## Example Session

Here's what a typical session looks like:

```bash
# 1. Connect to your cloud computer
$ lfr connect alice
Connecting to alice's instance in project cs101...
Instance state: running
Connecting to 35.87.81.251...
✅ Connected!

# 2. You're now on your cloud computer
alice@instance:~$ pwd
/home/alice

# 3. Create and edit a Python file
alice@instance:~$ nano hello.py
# Type your Python code, save with Ctrl+X, Y, Enter

# 4. Run your program
alice@instance:~$ python3 hello.py
Hello, World!

# 5. Check shared class files
alice@instance:~$ ls /mnt/efs/shared/
assignment1.pdf  class_notes.txt  dataset.csv

# 6. When done, exit
alice@instance:~$ exit
Connection to 35.87.81.251 closed.

# 7. You're back on your own computer
$
```

## Getting More Help

- **Basic Commands**: Ask your teacher for a "Linux commands cheat sheet"
- **Programming Help**: Use your favorite programming tutorials
- **Technical Issues**: Contact your teacher or TA
- **Advanced Features**: Read the detailed documentation

Remember: Everyone starts as a beginner. Don't be afraid to ask questions and experiment!

## Troubleshooting Quick Fixes

**"No access token found"**
- Run the activate command again with your token

**"Instance is not running"**
- Wait 1-2 minutes and try again
- Ask your teacher to start the class computers

**"Connection failed"**
- Check your internet connection
- Make sure you typed your username correctly
- Try again in a few minutes

**"Permission denied"**
- You might be trying to access someone else's files
- Ask your teacher about file permissions

Welcome to cloud computing! 🎓