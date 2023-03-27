import argparse
import os
import pathlib
import glob
import shutil
import subprocess


# Check that we are in bin folder
if pathlib.Path(os.getcwd()).name != "bin":
    print("Run the script inside the bin directory")
    exit(1)

# Configure argparse
parser = argparse.ArgumentParser(
    prog="Build",
    description="Build script for MSTK",
    epilog="Made by github.com/nfwGytautas/MSTK"
)

parser.add_argument("target", choices=["build", "clean"])

args = parser.parse_args()

# Copy the gomods directory into tmp
shutil.copytree("../gomods/", "./tmp/")

# Get all packages
gomods = pathlib.Path("../gomods/")
balancers = pathlib.Path("../gomods/balancers/")

packages = [x.name for x in gomods.glob("*") if not x.name.endswith("-api") and not x.name == "README.md" and not x.name == "balancers"]
balancers = ["balancers/" + x.name for x in balancers.glob("*") if not x.name == "README.md"]
packages = packages + balancers

print("Found", len(packages), "packages")

# Script content
if args.target == "build":
    print("Running build")

    for i, package in enumerate(packages):
        print(f"{i+1}/{len(packages)} {package}")
        os.system(f"go build -o ./{package} ./tmp/{package}/*.go")

elif args.target == "clean":
    print("Running clean")

    for i, package in enumerate(packages):
        print(f"{i+1}/{len(packages)} {package}")
        if pathlib.Path(package).exists():
            os.remove(package)

# Delete tmp directory
shutil.rmtree("./tmp/", ignore_errors=False, onerror=None)
shutil.rmtree("./balancers/", ignore_errors=True, onerror=None)
