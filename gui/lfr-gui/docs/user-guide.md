# LFR Tools Desktop App - User Guide

Welcome to the LFR Tools desktop application! This guide helps you use the visual interface to manage your cloud computers.

## Getting Started

### Installing the Desktop App

**On Mac:**
1. Download LFR Tools from the releases page
2. Open the `.dmg` file
3. Drag LFR Tools to your Applications folder
4. Open LFR Tools from Applications

**On Windows:**
1. Download the Windows installer
2. Run the installer (you may need to allow it in Windows Security)
3. Find LFR Tools in your Start menu
4. Click to open

**On Linux:**
1. Download the `.AppImage` file
2. Make it executable: `chmod +x LFR-Tools.AppImage`
3. Double-click to run, or run from terminal

### First Time Setup

When you first open LFR Tools:

1. **Choose your role:**
   - **Student**: If you received an access code from your teacher
   - **Teacher/Professor**: If you have AWS account access
   - **TA**: If your professor gave you helper access

2. **Enter your information:**
   - AWS profile (usually "aws")
   - Your username
   - Project/class name

3. **Test connection:**
   - The app will verify it can connect to AWS
   - You'll see a green checkmark if everything works

## For Students

### Your Dashboard

When you open LFR Tools as a student, you'll see:

#### **Main Screen**
```
┌─ My Cloud Computer ─────────────────────────────────┐
│                                                     │
│ 📚 CS101 Fall 2024               🟢 Ready to use   │
│                                                     │
│ 💻 alice-ubuntu-22-04                              │
│ Status: Ready to connect                            │
│ Budget: $18.50 / $25.00 (74% used)                 │
│ Time Left: 45 days                                  │
│                                                     │
│ [🔗 Connect Now] [📁 Shared Files] [❓ Get Help]   │
│                                                     │
│ 💡 Tips:                                           │
│ • Your computer sleeps automatically to save money │
│ • Shared files are in /mnt/efs/shared/            │
│ • Save your work often                             │
│                                                     │
└─────────────────────────────────────────────────────┘
```

#### **Connecting to Your Computer**

1. **Click "Connect Now"**
   - If your computer is sleeping, it will wake up (30-60 seconds)
   - You'll see a progress indicator

2. **Terminal Opens**
   - A black terminal window appears in the app
   - You'll see: `alice@instance:~$`
   - You can now type commands just like on your own computer

3. **Working in the Terminal**
   ```bash
   # See what's in your folder
   ls

   # Create a new Python file
   nano homework1.py

   # Run your Python program
   python3 homework1.py

   # Check shared class files
   ls /mnt/efs/shared/
   ```

#### **Understanding Your Budget**

Your budget bar shows:
- **Green**: You're doing great with money
- **Yellow**: Be careful, you're using a lot of budget
- **Red**: Almost out of money, talk to your teacher

**Tips to save money:**
- Close the terminal when you're done working
- Your computer sleeps automatically after 2 hours
- Don't leave your computer running overnight

### Common Student Tasks

#### **Doing Homework**
1. Click "Connect Now"
2. Wait for terminal to open
3. Create your homework file: `nano assignment1.py`
4. Write your code
5. Save with `Ctrl+X`, then `Y`, then `Enter`
6. Test your program: `python3 assignment1.py`
7. When done, type `exit` to disconnect

#### **Accessing Shared Files**
1. Connect to your computer
2. Look at shared files: `ls /mnt/efs/shared/`
3. Copy files you need: `cp /mnt/efs/shared/assignment.py ~/`
4. Work on the copied file
5. Submit your work: `cp my_homework.py /mnt/efs/submissions/`

#### **Getting Help**
1. Click "Get Help" button for quick tips
2. If your computer won't start, wait 2-3 minutes and try again
3. If you can't connect, ask your teacher or TA
4. If you lost your work, check `/mnt/efs/shared/` folder

## For Teachers and Professors

### Main Dashboard

As a teacher, you see a comprehensive class management interface:

#### **Class Overview**
```
┌─ CS101 Fall 2024 Overview ─────────────────────────┐
│                                                     │
│ 📊 Summary          💰 Budget          🚨 Alerts   │
│ Students: 25        Used: $340/$500    3 requests  │
│ Online: 8          Remaining: $160     1 budget    │
│ Running: 8         Days: 45           1 help       │
│                                                     │
│ [▶️ Start All] [⏹️ Stop All] [📊 Analytics]        │
│                                                     │
└─────────────────────────────────────────────────────┘
```

#### **Student Management Table**
- See all students and their computer status
- Start/stop individual or multiple computers
- Connect directly to help struggling students
- Monitor budget usage per student

#### **Real-Time Monitoring**
- See who's currently working
- Get alerts when students need help
- Monitor class activity during lab sessions
- Track budget usage in real-time

### Common Teacher Tasks

#### **Starting Class**
1. Open LFR Tools
2. Go to "Instances" tab
3. Click "Start All" button
4. Wait for all computers to start (1-2 minutes)
5. Tell students they can now connect

#### **Helping a Student**
1. Find the student in the Instances table
2. Click the "SSH" button next to their name
3. A terminal opens connected to their computer
4. You can see their files and help debug
5. Type `exit` when done helping

#### **Managing Costs**
1. Go to "Analytics" tab
2. See spending trends and projections
3. Click "Apply Optimization Suggestions"
4. Review and approve cost-saving changes
5. Monitor budget alerts and warnings

#### **End of Class**
1. Go to "Instances" tab
2. Click "Stop All" button
3. Computers turn off to save money
4. Students can still connect later (computers will auto-start)

### Advanced Features

#### **Bulk Operations**
1. Go to "Users & Groups" tab
2. Upload a CSV file with student information
3. Preview the changes before applying
4. Create multiple students at once

#### **Software Installation**
1. Go to "Software" tab (coming soon)
2. Select software pack (Python, Data Science, etc.)
3. Choose which students to install on
4. Monitor installation progress

#### **Cost Analytics**
1. Go to "Analytics" tab
2. See detailed cost breakdowns
3. Identify students using too many resources
4. Apply optimization suggestions automatically

## For TAs and Lab Assistants

### Support Dashboard

As a TA, you see a student-focused interface:

#### **Student Support Center**
```
┌─ Student Support ───────────────────────────────────┐
│                                                     │
│ 🆘 Help Requests                                   │
│ • alice: Can't run Python (2 min ago)             │
│ • bob: Software installation issue (5 min ago)     │
│                                                     │
│ 📊 Lab Session: "Python Functions"                │
│ Online: 22/25 students                             │
│ Missing: alice, bob, charlie                       │
│                                                     │
│ [📞 Help alice] [▶️ Start Missing] [📢 Message]    │
│                                                     │
└─────────────────────────────────────────────────────┘
```

#### **TA Capabilities**
- Start and stop student computers
- Connect to help students directly
- Monitor lab session attendance
- Cannot create or delete accounts
- Cannot modify budgets or settings

### Common TA Tasks

#### **Helping During Lab**
1. Watch the "Support Requests" section
2. Click "Help [student]" to connect to their computer
3. Debug their code or fix installation issues
4. Guide them through the solution
5. Close connection when problem is solved

#### **Managing Lab Sessions**
1. Before lab: Check that all computers are ready
2. During lab: Monitor attendance and help requests
3. Contact missing students if needed
4. Report any technical issues to the professor

## Tips for Everyone

### **Performance Tips**
- The app loads faster after the first time (cached)
- Close unused tabs to save computer memory
- Use "Refresh" buttons to get the latest information
- Keep the app updated for best performance

### **Troubleshooting**

#### **App Won't Start**
1. Make sure you have internet connection
2. Check that you have the latest version
3. Try restarting your computer
4. Contact technical support

#### **Can't See Your Project/Class**
1. Check that you selected the right role
2. Verify your AWS profile is correct
3. Ask your teacher if your access is set up correctly

#### **Terminal Won't Connect**
1. Check that your cloud computer is running (green status)
2. Wait 1-2 minutes if it's starting up
3. Try refreshing the app
4. Ask your teacher to check the instance status

#### **Slow Performance**
1. Close other applications using lots of memory
2. Check your internet connection speed
3. Try refreshing the app
4. Clear your browser cache (in app settings)

### **Security Notes**

- **Students**: Your access only works on your specific computer
- **Teachers**: Log out when using shared computers
- **Everyone**: Don't share your login information
- **Data**: All connections are encrypted and secure

### **Getting More Help**

- **In-App Help**: Click "Get Help" buttons throughout the app
- **Documentation**: Check the complete guides in the docs folder
- **Technical Support**: Contact your IT department for AWS issues
- **Bug Reports**: Report problems on the GitHub page

Remember: This app is designed to make cloud computing simple and accessible. Don't be afraid to explore and try different features! 🎓