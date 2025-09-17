# Getting Started with LFR Tools

Welcome! This guide will help you start using LFR Tools in simple steps.

## What is LFR Tools?

LFR Tools helps you manage cloud computers on Amazon Web Services (AWS). Think of it like having a remote computer that you can turn on and off, just like your laptop, but it lives in the cloud.

### Who Uses LFR Tools?

- **Teachers**: Set up cloud computers for students
- **Students**: Connect to your assigned cloud computer
- **Researchers**: Manage powerful computers for data analysis
- **IT Admins**: Control costs and manage many users

## Quick Start for Students

If your teacher gave you an access code, follow these steps:

### Step 1: Install LFR Tools

**On Mac:**
```bash
# Open Terminal and type:
brew install lfr
```

**On Windows:**
1. Download LFR Tools from the releases page
2. Run the installer
3. Open Command Prompt

**On Linux:**
```bash
# Download and install:
curl -L https://github.com/scttfrdmn/lfr-tools/releases/latest/download/lfr_Linux_x86_64.tar.gz | tar xz
sudo mv lfr /usr/local/bin/
```

### Step 2: Activate Your Access

Your teacher will give you two things:
1. **Access Token**: A long code like `cs101-alice-AbCdEf123`
2. **Student ID**: Your school ID number

```bash
# Type this in your terminal:
lfr connect activate <your-access-token> <your-student-id>

# Example:
lfr connect activate cs101-alice-AbCdEf123 12345
```

This connects the access code to your computer. You only need to do this once.

### Step 3: Connect to Your Cloud Computer

```bash
# Connect to your assigned computer:
lfr connect alice

# Replace "alice" with your username
```

**What happens:**
- If your computer is sleeping, it wakes up (takes 30-60 seconds)
- You get connected automatically
- You see a command line, just like Terminal

### Step 4: Use Your Cloud Computer

Once connected, you can:
- Run programs: `python3`, `node`, `git`
- Edit files: `nano myfile.py` or `vim myfile.py`
- Install software your teacher allows
- Save your work (it stays on the cloud computer)

### Step 5: Disconnect

To leave your cloud computer:
```bash
# Type:
exit

# Or just close the terminal window
```

Your work is saved automatically. The computer might go to sleep after a while to save money.

## Quick Start for Teachers

### Step 1: Set Up Your Class

```bash
# Create a class environment:
lfr students setup environment \
  --project=cs101-fall2024 \
  --s3-bucket=cs101-status \
  --students=alice,bob,charlie \
  --professor=drsmith

# This creates the foundation for student access
```

### Step 2: Create Student Accounts

```bash
# Create a CSV file with student info:
lfr users template students.csv

# Edit students.csv with your student names, then:
lfr users create-bulk students.csv --start-stopped

# --start-stopped saves money by keeping computers off until needed
```

### Step 3: Generate Student Access Codes

```bash
# Create secure access tokens:
lfr students generate tokens --project=cs101-fall2024

# This creates a file with access codes for each student
```

### Step 4: Distribute Access to Students

1. Send each student their access token via email
2. Include these instructions:
   - Install LFR Tools
   - Run: `lfr connect activate <their-token> <their-student-id>`
   - Connect with: `lfr connect <their-username>`

### Step 5: Manage Your Class

```bash
# Start computers for class:
lfr instances start --project=cs101-fall2024 --wait

# Check who's requesting access:
lfr students check requests --project=cs101-fall2024

# See class status:
lfr students status --project=cs101-fall2024

# Stop computers after class (saves money):
lfr instances stop --project=cs101-fall2024
```

## Basic Concepts

### Cloud Computers (Instances)
Think of these as remote computers you can rent by the hour. They have:
- **CPU**: How fast they think (like your laptop's processor)
- **Memory (RAM)**: How much they can remember at once
- **Storage**: Where files are saved
- **Operating System**: Usually Ubuntu Linux

### Projects
A way to group related computers. Examples:
- `cs101-fall2024` - A computer science class
- `research-lab` - A research project
- `web-development` - A development team

### Users and Permissions
- **Professors**: Can create and delete everything
- **TAs**: Can help students and start/stop computers
- **Students**: Can only connect to their assigned computer

### Storage Types
- **Instance Storage**: Saved on your specific computer (goes away if computer is deleted)
- **EFS Shared Storage**: Shared between all class members (like Google Drive)
- **EBS Volumes**: Extra storage you can attach (like plugging in a USB drive)

## Common Tasks

### For Students

**Connect to Your Computer:**
```bash
lfr connect alice
```

**Check Available Connections:**
```bash
lfr connect list
```

**Get Help:**
```bash
lfr connect --help
```

### For Teachers

**See All Students:**
```bash
lfr users list --project=cs101
```

**Start Computers for Class:**
```bash
lfr instances start --project=cs101 --wait
```

**Install Software on All Computers:**
```bash
lfr software install python-dev alice
lfr software install python-dev bob
# Or create a script to install on everyone
```

**Check Costs:**
```bash
lfr idle advanced analyze --project=cs101
```

## Safety Features

### For Students
- You can only access your assigned computer
- Your access expires at the end of the semester
- Your access token only works on your specific device

### For Teachers
- Students cannot create or delete computers
- All actions are logged
- Budget limits prevent surprise costs
- Computers automatically turn off when not used

## Getting Help

If something doesn't work:

1. **Check the error message** - it usually tells you what's wrong
2. **Make sure you're connected to the internet**
3. **Try the command again** - sometimes it's a temporary problem
4. **Ask your teacher or TA** - they can check your computer's status
5. **Check the troubleshooting guide** - common problems and solutions

## Important Notes

### For Students
- **Save your work often** - computers can turn off automatically
- **Don't share your access token** - it only works on your device
- **Ask for help early** - teachers and TAs are there to help
- **Use the shared folder** for group projects

### For Teachers
- **Set budgets early** - prevent surprise costs
- **Test everything first** - try all steps before class
- **Have backup plans** - technology sometimes breaks
- **Monitor usage** - help students who are struggling

## Next Steps

- **Students**: Try connecting and exploring your cloud computer
- **Teachers**: Set up a test class with a few sample students
- **Everyone**: Read the detailed tutorials for your specific needs

Welcome to the world of cloud computing! 🚀