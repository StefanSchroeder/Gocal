{ buildGoModule
, fetchFromGitHub
, lib
}:
buildGoModule rec {
	pname = "gocal";
	version = "0.9.0";

	src = fetchFromGitHub {
		owner = "StefanSchroeder";
		repo = "Gocal";
		rev = "v${version}";
		sha256 = "sha256-hTAZhpVfQSgXwPY6K9mDmJFdQcBRRhvW8TnntjrZrrE=";

	};
	vendorSha256 = "sha256-pNt7goaWR6gdf25t/SqwRQ+MYNdqOaUhRLKriAr2ocE=";
	meta = with lib; {
		description = "Create PDF calendars";
		homepage = "https://github.com/StefanSchroeder/Gocal";
		license = licenses.mit;
		maintainers = with maintainers; [ stefan ];
		mainProgram = "gocalendar";
	};
}
