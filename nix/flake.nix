{
	description = "My nix-darwin system flake";

	inputs = {
		# Main nix
		nixpkgs.follows = "pkgs-unstable";
		pkgs-stable.url = "github:nixos/nixpkgs/nixos-24.05";
		pkgs-unstable.url = "github:nixos/nixpkgs/nixos-unstable";

		# Other
		disko.url = "github:nix-community/disko";
		disko.inputs.nixpkgs.follows = "nixpkgs";
	};

	outputs = inputs: {
		# nixos-anywhere --flake .#nixos-boreas root@<ip>
		nixosConfigurations."nixos-boreas" = inputs.nixpkgs.lib.nixosSystem {
			system = "x86_64-linux";
			modules = [
				inputs.disko.nixosModules.disko
				./configuration.nix
			];
		};
	};
}
