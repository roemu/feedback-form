{ modulesPath, lib, pkgs, ... }:
{
	imports = [
		(modulesPath + "/installer/scan/not-detected.nix")
		(modulesPath + "/profiles/qemu-guest.nix")
		./disk-config.nix
	];
	nix.settings.extra-experimental-features = "flakes nix-command";
	boot.loader.grub = {
		efiSupport = true;
		efiInstallAsRemovable = true;
	};
	services.openssh = {
		enable = true;
		ports = [2222];
	};

	environment.systemPackages = map lib.lowPrio [
		pkgs.curl
		pkgs.gitMinimal
		pkgs.sqlite
		pkgs.go
	];

	programs.neovim = {
		enable = true;
		defaultEditor = true;
		viAlias = true;
	};

	users.users.root = {
		openssh.authorizedKeys.keys = [
			"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOIcrTOdoPmASCfBPjt+qm/iGQ6ASExs1YtOAtKIMJty"
			"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOzG8TYzMry09a1s6IAqO+N3+tSKaAcvwqB7i2vWCpT7"
			"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICKnGNGP2GQeA1FCp7u3OccD8u6hQXFqJZW9rd0GJcZe"
		];
		hashedPassword = "$y$j9T$JLjEK1XtXY8qaDflocaE5.$zOhlVgvtbbPo3uGeMiSWv4N/EAP/NrCmKYmsNqGBF33";
	};
	users.users.roemu = {
		group = "roemu";
		isNormalUser = true;
		openssh.authorizedKeys.keys = [
			"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOIcrTOdoPmASCfBPjt+qm/iGQ6ASExs1YtOAtKIMJty"
			"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOzG8TYzMry09a1s6IAqO+N3+tSKaAcvwqB7i2vWCpT7"
			"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICKnGNGP2GQeA1FCp7u3OccD8u6hQXFqJZW9rd0GJcZe"
		];
		extraGroups = ["wheel"];
		hashedPassword = "$y$j9T$JLjEK1XtXY8qaDflocaE5.$zOhlVgvtbbPo3uGeMiSWv4N/EAP/NrCmKYmsNqGBF33";
	};
	users.groups.roemu = {};

	networking.hostName = "nixos-boreas";
	networking.firewall = {
		enable = true;
		allowedTCPPorts = [22 2222 443 80];
	};

	system.stateVersion = "24.05";
}
