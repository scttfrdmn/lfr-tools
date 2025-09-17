# Teacher Guide: Managing Your Digital Classroom

This guide helps teachers set up and manage cloud computers for their students.

## What You Can Do

As a teacher, you can:
- **Create cloud computers** for all your students
- **Install software** that students need for class
- **Control costs** by managing when computers run
- **Share files** with the whole class
- **Help students** when they have problems

## Setting Up Your First Class

### Before You Start

Make sure you have:
- AWS account access (ask your IT department)
- List of student usernames
- An idea of what software students need

### Step 1: Install LFR Tools

```bash
# On Mac:
brew install lfr

# Follow the installation guide for other platforms
```

### Step 2: Set Up Your Class Environment

```bash
# Create your class:
lfr students setup environment \
  --project=cs101-fall2024 \
  --s3-bucket=cs101-status-bucket \
  --students=alice,bob,charlie,david \
  --start-date=2024-08-15 \
  --end-date=2024-12-15
```

**What this does:**
- Creates a secure communication system for students
- Sets up automatic cost controls
- Prepares the foundation for student access

### Step 3: Create Student Computers

```bash
# Create a template for your students:
lfr users template students.csv

# Edit the file with your student information, then:
lfr users create-bulk students.csv --start-stopped
```

**The `--start-stopped` flag is important**: It creates the computers but keeps them turned off to save money. Students can turn them on when needed.

### Step 4: Install Software

```bash
# See available software packs:
lfr software list

# Install Python environment on Alice's computer:
lfr software install python-dev alice

# Install for multiple students:
lfr software install python-dev alice
lfr software install python-dev bob
lfr software install python-dev charlie
```

**Available software packs:**
- `python-dev`: Python, pip, common libraries
- `data-science`: Python + R + Jupyter + data tools
- `web-dev`: Node.js, npm, web development tools
- `gpu-ml`: Machine learning tools for GPU computers

### Step 5: Generate Student Access

```bash
# Create secure access tokens:
lfr students generate tokens --project=cs101-fall2024

# This creates a file with tokens like:
# alice:student:cs101-alice-AbCdEf123
# bob:student:cs101-bob-XyZ789def
```

### Step 6: Give Students Access

Send each student their specific token via email:

**Email template:**
```
Subject: Your CS101 Cloud Computer Access

Hi Alice,

Here's how to connect to your cloud computer for CS101:

1. Install LFR Tools: brew install lfr
2. Set up access: lfr connect activate cs101-alice-AbCdEf123 <your-student-id>
3. Connect anytime: lfr connect alice

Your computer will automatically start when you connect.

Let me know if you have any problems!

Professor Smith
```

## Managing Your Class

### Daily Class Management

**Before class starts:**
```bash
# Start all student computers:
lfr instances start --project=cs101-fall2024 --wait

# This takes about 1-2 minutes for all computers to start
```

**During class:**
```bash
# Check who's online:
lfr students status --project=cs101-fall2024

# Help a specific student:
lfr ssh connect alice  # Connect to Alice's computer to help

# Check if someone requested access:
lfr students check requests --project=cs101-fall2024 --auto-approve
```

**After class:**
```bash
# Turn off computers to save money:
lfr instances stop --project=cs101-fall2024
```

### Managing Files and Storage

**Set up shared storage:**
```bash
# Create shared folder for class materials:
lfr efs create class-materials --project=cs101-fall2024

# Mount on all student computers:
lfr efs mount-all fs-12345 --project=cs101-fall2024
```

**File organization suggestions:**
```
/mnt/efs/shared/
├── lectures/          # Your lecture materials (read-only for students)
├── assignments/       # Assignment descriptions
├── examples/          # Code examples
└── submissions/       # Where students submit work
```

### Cost Management

**Check spending:**
```bash
# See current costs:
lfr idle advanced analyze --project=cs101-fall2024

# Set up automatic cost controls:
lfr idle advanced policies apply educational-conservative --project=cs101-fall2024
```

**Cost-saving tips:**
- Use `--start-stopped` when creating student computers
- Turn off computers after class
- Use shared storage instead of individual storage when possible
- Consider smaller computer sizes for basic programming classes

### Adding New Students

**Mid-semester additions:**
```bash
# Add new student to existing class:
lfr users create --project=cs101-fall2024 \
  --blueprint=ubuntu_22_04 \
  --bundle=app_standard_xl_1_0 \
  --region=us-west-2 \
  --users=new_student

# Generate their access token:
lfr students generate tokens --project=cs101-fall2024
```

### Handling Problems

**Student can't connect:**
1. Check if their computer is running: `lfr instances list --user=alice`
2. Start their computer: `lfr instances start --users=alice`
3. Check their access token hasn't expired

**Student lost their work:**
1. Connect to their computer: `lfr ssh connect alice`
2. Look for their files: `find /home/alice -name "*.py"`
3. Check shared folders: `ls /mnt/efs/shared/`

**Too expensive:**
1. Analyze usage: `lfr idle advanced analyze --project=cs101`
2. Apply cost-saving policies: `lfr idle advanced policies apply educational-conservative`
3. Resize oversized computers: `lfr instances resize expensive-computer down`

## Advanced Features

### Working with TAs

You can give TAs limited access to help manage the class:

```bash
# When generating tokens, include TAs:
lfr students setup environment \
  --project=cs101-fall2024 \
  --tas=ta-alice,ta-bob

# TAs can then:
# - Start and stop student computers
# - Check student status
# - Help with technical issues
# - Cannot create or delete accounts
```

### Semester Management

**End of semester cleanup:**
```bash
# Create final snapshots (backup student work):
lfr instances snapshot alice-ubuntu
lfr instances snapshot bob-ubuntu

# Generate final cost report:
lfr idle advanced analyze --project=cs101-fall2024 --days=120

# Clean up (when semester ends):
lfr users remove-bulk --project=cs101-fall2024 --confirm
```

### Different Computer Types

**For basic programming:**
- `app_standard_xl_1_0` - Good for Python, web development
- Cost: ~$25/month if used 8 hours/day

**For data science:**
- `app_standard_2xl_1_0` - More memory for data analysis
- Cost: ~$40/month if used 8 hours/day

**For machine learning:**
- `gpu_nvidia_xl_1_0` - Has GPU for AI/ML work
- Cost: ~$80/month if used 8 hours/day

### Assignment Workflows

**Distributing assignments:**
1. Put assignment files in `/mnt/efs/shared/assignments/`
2. Students copy to their own folders
3. Students work on their own computers
4. Students submit to `/mnt/efs/submissions/`

**Collecting submissions:**
1. Check submissions folder
2. Download or review work
3. Provide feedback via shared folder

## Security and Safety

### What Students Cannot Do
- Access other students' computers
- Create or delete computers
- Change class settings
- Access your teacher account
- Share their access with others (tokens are locked to their device)

### What You Control
- When computers start and stop
- What software is installed
- How much money can be spent
- Who has access to what files
- When the class access expires

### Best Practices
- Always test new software before class
- Set spending limits to prevent surprises
- Turn off computers when not in use
- Keep backup copies of important class materials
- Monitor student usage patterns

## Getting Help

### For Technical Issues
- Check the troubleshooting guide
- Use `lfr --debug` commands for detailed error information
- Contact your IT department for AWS account issues

### For Teaching Issues
- Start small with a few test students
- Have a backup plan for technology failures
- Prepare offline alternatives for important lessons

### Community Support
- Check the documentation website
- Ask questions in the GitHub discussions
- Share your experiences with other educators

Remember: You're giving students access to professional development tools. This prepares them for real-world software development! 🎓