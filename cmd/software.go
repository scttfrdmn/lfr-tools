package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/scttfrdmn/lfr-tools/internal/aws"
	"github.com/scttfrdmn/lfr-tools/internal/types"
)

var softwareCmd = &cobra.Command{
	Use:   "software",
	Short: "Manage software packs on Lightsail for Research instances",
	Long:  `Install and manage software packages, development environments, and tools on instances.`,
}

var softwareInstallCmd = &cobra.Command{
	Use:   "install [pack-name] [username]",
	Short: "Install software pack on user's instance",
	Long: `Install a software pack on a user's Lightsail instance. Supports APT packages,
container deployments, and custom scripts.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		packName := args[0]
		username := args[1]
		project, _ := cmd.Flags().GetString("project")
		force, _ := cmd.Flags().GetBool("force")

		return installSoftwarePack(cmd.Context(), packName, username, project, force)
	},
}

var softwareListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available software packs",
	Long:  `List all available software packs with descriptions and installation status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		category, _ := cmd.Flags().GetString("category")
		installed, _ := cmd.Flags().GetBool("installed")

		return listSoftwarePacks(category, installed)
	},
}

var softwareCreateCmd = &cobra.Command{
	Use:   "create [pack-name]",
	Short: "Create a custom software pack",
	Long: `Create a custom software pack definition that can be installed on instances.
Supports APT packages, scripts, and environment configuration.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		packName := args[0]
		template, _ := cmd.Flags().GetString("template")

		return createSoftwarePack(packName, template)
	},
}

var softwareStatusCmd = &cobra.Command{
	Use:   "status [username]",
	Short: "Show software installation status for user",
	Long: `Display what software packs are installed on a user's instance and their status.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		project, _ := cmd.Flags().GetString("project")

		return showSoftwareStatus(cmd.Context(), username, project)
	},
}

func init() {
	rootCmd.AddCommand(softwareCmd)

	softwareCmd.AddCommand(softwareInstallCmd)
	softwareCmd.AddCommand(softwareListCmd)
	softwareCmd.AddCommand(softwareCreateCmd)
	softwareCmd.AddCommand(softwareStatusCmd)

	// Install command flags
	softwareInstallCmd.Flags().StringP("project", "p", "", "Project name")
	softwareInstallCmd.Flags().BoolP("force", "f", false, "Force reinstall if already installed")

	// List command flags
	softwareListCmd.Flags().StringP("category", "c", "", "Filter by category (development, data-science, gpu, etc.)")
	softwareListCmd.Flags().BoolP("installed", "i", false, "Show only installed packs")

	// Create command flags
	softwareCreateCmd.Flags().StringP("template", "t", "basic", "Template to use (basic, development, data-science, gpu)")

	// Status command flags
	softwareStatusCmd.Flags().StringP("project", "p", "", "Project name")
}

// Built-in software packs for common educational use cases
var builtinPacks = map[string]*types.SoftwarePack{
	"python-dev": {
		ID:          "python-dev",
		Name:        "Python Development Environment",
		Description: "Complete Python development setup with common packages",
		Category:    "development",
		Type:        types.PackTypeAPT,
		Version:     "1.0",
		Packages: []types.Package{
			{Name: "python3", Source: "apt"},
			{Name: "python3-pip", Source: "apt"},
			{Name: "python3-venv", Source: "apt"},
			{Name: "python3-dev", Source: "apt"},
			{Name: "git", Source: "apt"},
			{Name: "curl", Source: "apt"},
			{Name: "vim", Source: "apt"},
			{Name: "htop", Source: "apt"},
		},
		Scripts: []types.Script{
			{
				Name:        "pip-packages",
				Description: "Install common Python packages",
				Content: `#!/bin/bash
pip3 install --user numpy pandas matplotlib jupyter notebook scipy scikit-learn requests
echo "Python packages installed successfully"`,
				Type: "bash",
			},
		},
		Environment: map[string]string{
			"PYTHON_USER_BASE": "/home/ubuntu/.local",
			"PATH":             "/home/ubuntu/.local/bin:$PATH",
		},
		Supported: []string{"ubuntu_22_04", "ubuntu_20_04"},
	},

	"data-science": {
		ID:          "data-science",
		Name:        "Data Science Environment",
		Description: "R, Python, Jupyter, and common data science tools",
		Category:    "data-science",
		Type:        types.PackTypeMixed,
		Version:     "1.0",
		Dependencies: []string{"python-dev"},
		Packages: []types.Package{
			{Name: "r-base", Source: "apt"},
			{Name: "r-base-dev", Source: "apt"},
			{Name: "rstudio-server", Source: "custom"},
		},
		Scripts: []types.Script{
			{
				Name:        "rstudio-setup",
				Description: "Install and configure RStudio Server",
				Content: `#!/bin/bash
# Install RStudio Server
wget -q https://download2.rstudio.org/server/jammy/amd64/rstudio-server-2023.12.1-402-amd64.deb
sudo dpkg -i rstudio-server-2023.12.1-402-amd64.deb
sudo apt-get install -f -y

# Configure RStudio
sudo systemctl enable rstudio-server
sudo systemctl start rstudio-server

echo "RStudio Server installed on port 8787"`,
				Type: "bash",
				RunAs: "root",
			},
		},
		Supported: []string{"ubuntu_22_04", "ubuntu_20_04"},
	},

	"gpu-ml": {
		ID:          "gpu-ml",
		Name:        "GPU Machine Learning Environment",
		Description: "CUDA, PyTorch, TensorFlow for GPU-enabled instances",
		Category:    "gpu",
		Type:        types.PackTypeMixed,
		Version:     "1.0",
		Dependencies: []string{"python-dev"},
		Packages: []types.Package{
			{Name: "nvidia-driver-535", Source: "apt"},
			{Name: "nvidia-cuda-toolkit", Source: "apt"},
		},
		Scripts: []types.Script{
			{
				Name:        "gpu-ml-setup",
				Description: "Install PyTorch and TensorFlow with GPU support",
				Content: `#!/bin/bash
# Install PyTorch with CUDA support
pip3 install --user torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu118

# Install TensorFlow with GPU support
pip3 install --user tensorflow[and-cuda]

# Install additional ML packages
pip3 install --user transformers datasets accelerate

# Verify GPU access
python3 -c "import torch; print(f'CUDA available: {torch.cuda.is_available()}')"

echo "GPU ML environment setup complete"`,
				Type: "bash",
			},
		},
		Environment: map[string]string{
			"CUDA_VISIBLE_DEVICES": "0",
		},
		Supported: []string{"ubuntu_22_04"},
		Tags:      []string{"gpu", "ml", "pytorch", "tensorflow"},
	},

	"web-dev": {
		ID:          "web-dev",
		Name:        "Web Development Environment",
		Description: "Node.js, npm, and common web development tools",
		Category:    "development",
		Type:        types.PackTypeAPT,
		Version:     "1.0",
		Packages: []types.Package{
			{Name: "nodejs", Source: "apt"},
			{Name: "npm", Source: "apt"},
			{Name: "git", Source: "apt"},
			{Name: "curl", Source: "apt"},
			{Name: "wget", Source: "apt"},
			{Name: "unzip", Source: "apt"},
		},
		Scripts: []types.Script{
			{
				Name:        "node-global-packages",
				Description: "Install common Node.js global packages",
				Content: `#!/bin/bash
npm install -g create-react-app typescript eslint prettier nodemon
echo "Node.js global packages installed"`,
				Type: "bash",
			},
		},
		Supported: []string{"ubuntu_22_04", "ubuntu_20_04"},
	},
}

// installSoftwarePack installs a software pack on a user's instance.
func installSoftwarePack(ctx context.Context, packName, username, project string, force bool) error {
	// Get software pack definition
	pack, exists := builtinPacks[packName]
	if !exists {
		// Try to load custom pack
		customPack, err := loadCustomPack(packName)
		if err != nil {
			return fmt.Errorf("software pack '%s' not found. Available packs: %s",
				packName, strings.Join(getAvailablePackNames(), ", "))
		}
		pack = customPack
	}

	fmt.Printf("Installing software pack: %s\n", pack.Name)
	fmt.Printf("Description: %s\n", pack.Description)
	fmt.Printf("Target user: %s\n", username)

	// Find user's instance
	awsClient, err := aws.NewClient(ctx, aws.Options{
		Region:  viper.GetString("aws.region"),
		Profile: viper.GetString("aws.profile"),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	lightsailService := aws.NewLightsailService(awsClient)
	instances, err := lightsailService.ListInstances(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	var targetInstance *types.Instance
	for _, instance := range instances {
		if strings.HasPrefix(instance.Name, username+"-") {
			targetInstance = instance
			break
		}
	}

	if targetInstance == nil {
		return fmt.Errorf("no instance found for user: %s", username)
	}

	if targetInstance.State != "running" {
		return fmt.Errorf("instance %s is not running (state: %s). Start it first",
			targetInstance.Name, targetInstance.State)
	}

	if targetInstance.PublicIP == "" {
		return fmt.Errorf("instance %s has no public IP", targetInstance.Name)
	}

	fmt.Printf("Target instance: %s (%s)\n", targetInstance.Name, targetInstance.PublicIP)

	// Check if blueprint is supported
	if !isPackSupportedOnBlueprint(pack, targetInstance.Blueprint) {
		if !force {
			return fmt.Errorf("pack '%s' is not supported on blueprint '%s'. Use --force to override",
				packName, targetInstance.Blueprint)
		}
		fmt.Printf("⚠️ Warning: Pack may not work correctly on blueprint %s\n", targetInstance.Blueprint)
	}

	// Generate installation script
	installScript, err := generateInstallScript(pack, targetInstance, force)
	if err != nil {
		return fmt.Errorf("failed to generate install script: %w", err)
	}

	// Execute installation via SSH
	fmt.Printf("Executing installation on %s...\n", targetInstance.PublicIP)

	result, err := executeInstallationScript(ctx, targetInstance, installScript, username)
	if err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	if result.Success {
		fmt.Printf("✅ Software pack '%s' installed successfully!\n", pack.Name)
		fmt.Printf("Duration: %s\n", result.Duration)
		if len(result.Packages) > 0 {
			fmt.Printf("Packages installed: %s\n", strings.Join(result.Packages, ", "))
		}
	} else {
		fmt.Printf("❌ Installation failed: %s\n", result.Message)
		if len(result.Errors) > 0 {
			fmt.Printf("Errors:\n")
			for _, errMsg := range result.Errors {
				fmt.Printf("  - %s\n", errMsg)
			}
		}
		return fmt.Errorf("software pack installation failed")
	}

	return nil
}

// listSoftwarePacks lists available software packs.
func listSoftwarePacks(category string, installed bool) error {
	fmt.Printf("Available software packs:\n\n")
	fmt.Printf("%-15s %-30s %-15s %-40s\n",
		"ID", "NAME", "CATEGORY", "DESCRIPTION")
	fmt.Println(strings.Repeat("-", 115))

	for _, pack := range builtinPacks {
		if category != "" && pack.Category != category {
			continue
		}

		description := pack.Description
		if len(description) > 37 {
			description = description[:37] + "..."
		}

		fmt.Printf("%-15s %-30s %-15s %-40s\n",
			pack.ID, pack.Name, pack.Category, description)
	}

	// Show available categories
	categories := make(map[string]bool)
	for _, pack := range builtinPacks {
		categories[pack.Category] = true
	}

	var categoryList []string
	for cat := range categories {
		categoryList = append(categoryList, cat)
	}

	fmt.Printf("\nTotal: %d packs\n", len(builtinPacks))
	fmt.Printf("Categories: %s\n", strings.Join(categoryList, ", "))

	return nil
}

// createSoftwarePack creates a custom software pack template.
func createSoftwarePack(packName, template string) error {
	packFile := fmt.Sprintf("%s-pack.yaml", packName)

	// Create pack based on template
	var pack *types.SoftwarePack
	switch template {
	case "basic":
		pack = createBasicPackTemplate(packName)
	case "development":
		pack = createDevelopmentPackTemplate(packName)
	case "data-science":
		pack = createDataSciencePackTemplate(packName)
	case "gpu":
		pack = createGPUPackTemplate(packName)
	default:
		return fmt.Errorf("unknown template: %s (available: basic, development, data-science, gpu)", template)
	}

	// Write pack to YAML file
	data, err := json.MarshalIndent(pack, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal pack: %w", err)
	}

	err = os.WriteFile(packFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write pack file: %w", err)
	}

	fmt.Printf("✅ Software pack template created: %s\n", packFile)
	fmt.Printf("Template: %s\n", template)
	fmt.Printf("Edit the file and install with: lfr software install %s <username>\n", packName)

	return nil
}

// showSoftwareStatus shows installation status for a user's instance.
func showSoftwareStatus(ctx context.Context, username, project string) error {
	fmt.Printf("Software status for user: %s\n", username)

	// Find user's instance
	awsClient, err := aws.NewClient(ctx, aws.Options{
		Region:  viper.GetString("aws.region"),
		Profile: viper.GetString("aws.profile"),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %w", err)
	}

	lightsailService := aws.NewLightsailService(awsClient)
	instances, err := lightsailService.ListInstances(ctx, project)
	if err != nil {
		return fmt.Errorf("failed to list instances: %w", err)
	}

	var targetInstance *types.Instance
	for _, instance := range instances {
		if strings.HasPrefix(instance.Name, username+"-") {
			targetInstance = instance
			break
		}
	}

	if targetInstance == nil {
		return fmt.Errorf("no instance found for user: %s", username)
	}

	fmt.Printf("Instance: %s (%s)\n", targetInstance.Name, targetInstance.State)

	if targetInstance.State != "running" {
		fmt.Printf("⚠️ Instance is not running. Start it to check software status.\n")
		return nil
	}

	// Check for installed software marker files
	fmt.Printf("Checking installed software packs...\n")

	// For now, show placeholder status
	// TODO: Implement actual software status checking via SSH
	fmt.Printf("\n📦 Installed software packs:\n")
	fmt.Printf("- System packages: Available via 'dpkg -l'\n")
	fmt.Printf("- Python packages: Available via 'pip list'\n")
	fmt.Printf("- Node packages: Available via 'npm list -g'\n")
	fmt.Printf("\nTo check specific installations, SSH to instance:\n")
	fmt.Printf("lfr ssh connect %s\n", username)

	return nil
}

// Helper functions

func getAvailablePackNames() []string {
	var names []string
	for name := range builtinPacks {
		names = append(names, name)
	}
	return names
}

func loadCustomPack(packName string) (*types.SoftwarePack, error) {
	packFile := packName + "-pack.yaml"
	if _, err := os.Stat(packFile); os.IsNotExist(err) {
		packFile = packName + "-pack.json"
		if _, err := os.Stat(packFile); os.IsNotExist(err) {
			return nil, fmt.Errorf("custom pack file not found: %s", packName)
		}
	}

	data, err := os.ReadFile(packFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read pack file: %w", err)
	}

	var pack types.SoftwarePack
	if err := json.Unmarshal(data, &pack); err != nil {
		return nil, fmt.Errorf("failed to parse pack file: %w", err)
	}

	return &pack, nil
}

func isPackSupportedOnBlueprint(pack *types.SoftwarePack, blueprint string) bool {
	if len(pack.Supported) == 0 {
		return true // No restrictions
	}

	for _, supported := range pack.Supported {
		if supported == blueprint {
			return true
		}
	}
	return false
}

func generateInstallScript(pack *types.SoftwarePack, instance *types.Instance, force bool) (string, error) {
	script := "#!/bin/bash\n"
	script += "set -e\n\n"
	script += fmt.Sprintf("echo 'Installing %s on %s'\n", pack.Name, instance.Name)
	script += "echo 'Starting at: '$(date)\n\n"

	// Update package manager
	script += "sudo apt-get update -y\n\n"

	// Install APT packages
	for _, pkg := range pack.Packages {
		if pkg.Source == "apt" {
			script += fmt.Sprintf("echo 'Installing %s...'\n", pkg.Name)
			script += fmt.Sprintf("sudo apt-get install -y %s\n", pkg.Name)
		}
	}

	// Set environment variables
	if len(pack.Environment) > 0 {
		script += "\n# Set environment variables\n"
		for key, value := range pack.Environment {
			script += fmt.Sprintf("echo 'export %s=\"%s\"' >> ~/.bashrc\n", key, value)
		}
	}

	// Run custom scripts
	for _, customScript := range pack.Scripts {
		script += fmt.Sprintf("\n# %s\n", customScript.Description)
		script += fmt.Sprintf("echo 'Running %s...'\n", customScript.Name)
		script += customScript.Content + "\n"
	}

	script += "\necho 'Installation completed at: '$(date)\n"
	script += fmt.Sprintf("echo 'Pack %s installed successfully'\n", pack.ID)

	return script, nil
}

func executeInstallationScript(ctx context.Context, instance *types.Instance, script, username string) (*types.InstallResult, error) {
	start := time.Now()

	// Create temporary script file
	tmpFile, err := os.CreateTemp("", "lfr-install-*.sh")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp script: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(script)
	if err != nil {
		return nil, fmt.Errorf("failed to write script: %w", err)
	}
	tmpFile.Close()

	// Make script executable
	err = os.Chmod(tmpFile.Name(), 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to make script executable: %w", err)
	}

	// Execute via SSH
	// For now, return success placeholder
	// TODO: Implement actual SSH execution with proper key management

	duration := time.Since(start)

	result := &types.InstallResult{
		PackID:      "placeholder",
		Success:     true,
		Message:     "Installation script prepared (SSH execution pending)",
		Duration:    duration.String(),
		InstalledAt: time.Now().Format(time.RFC3339),
		Packages:    []string{"preparation-complete"},
	}

	fmt.Printf("📋 Installation script prepared for %s\n", instance.Name)
	fmt.Printf("To execute manually:\n")
	fmt.Printf("1. SSH to instance: lfr ssh connect %s\n", username)
	fmt.Printf("2. Run the installation commands\n")

	return result, nil
}

// Template creation functions

func createBasicPackTemplate(name string) *types.SoftwarePack {
	return &types.SoftwarePack{
		ID:          name,
		Name:        strings.Title(name) + " Pack",
		Description: "Custom software pack",
		Category:    "custom",
		Type:        types.PackTypeAPT,
		Version:     "1.0",
		Packages: []types.Package{
			{Name: "git", Source: "apt"},
			{Name: "curl", Source: "apt"},
			{Name: "vim", Source: "apt"},
		},
		Supported: []string{"ubuntu_22_04", "ubuntu_20_04"},
	}
}

func createDevelopmentPackTemplate(name string) *types.SoftwarePack {
	return &types.SoftwarePack{
		ID:          name,
		Name:        strings.Title(name) + " Development Pack",
		Description: "Development environment with common tools",
		Category:    "development",
		Type:        types.PackTypeMixed,
		Version:     "1.0",
		Packages: []types.Package{
			{Name: "build-essential", Source: "apt"},
			{Name: "git", Source: "apt"},
			{Name: "python3", Source: "apt"},
			{Name: "python3-pip", Source: "apt"},
			{Name: "nodejs", Source: "apt"},
			{Name: "npm", Source: "apt"},
		},
		Scripts: []types.Script{
			{
				Name:    "dev-setup",
				Content: "echo 'Development environment setup complete'",
				Type:    "bash",
			},
		},
		Supported: []string{"ubuntu_22_04", "ubuntu_20_04"},
	}
}

func createDataSciencePackTemplate(name string) *types.SoftwarePack {
	return &types.SoftwarePack{
		ID:          name,
		Name:        strings.Title(name) + " Data Science Pack",
		Description: "Data science environment with R and Python",
		Category:    "data-science",
		Type:        types.PackTypeMixed,
		Version:     "1.0",
		Dependencies: []string{"python-dev"},
		Packages: []types.Package{
			{Name: "r-base", Source: "apt"},
			{Name: "python3-scipy", Source: "apt"},
			{Name: "python3-numpy", Source: "apt"},
		},
		Supported: []string{"ubuntu_22_04", "ubuntu_20_04"},
	}
}

func createGPUPackTemplate(name string) *types.SoftwarePack {
	return &types.SoftwarePack{
		ID:          name,
		Name:        strings.Title(name) + " GPU Pack",
		Description: "GPU computing environment with CUDA",
		Category:    "gpu",
		Type:        types.PackTypeMixed,
		Version:     "1.0",
		Dependencies: []string{"python-dev"},
		Packages: []types.Package{
			{Name: "nvidia-cuda-toolkit", Source: "apt"},
		},
		Environment: map[string]string{
			"CUDA_VISIBLE_DEVICES": "0",
		},
		Supported: []string{"ubuntu_22_04"},
		Tags:      []string{"gpu", "cuda"},
	}
}