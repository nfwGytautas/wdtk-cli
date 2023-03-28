import argparse
import os
import pathlib
import glob
import shutil
import subprocess
import time
from string import Template
from datetime import datetime


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

parser.add_argument("target", choices=["build", "clean", "push"])

args = parser.parse_args()

start_time = time.time()

# Get all packages
gomods = pathlib.Path("../gomods/")
balancers = pathlib.Path("../gomods/balancers/")

packages = [x.name for x in gomods.glob("*") if not x.name.endswith("-api") and not x.name == "README.md" and not x.name == "balancers"]
balancers = [x.name for x in balancers.glob("*") if not x.name == "README.md"]

print(f"Found {len(packages)} packages, {len(balancers)} balancers")

# Script content
if args.target == "build":

    def build_target(idx, length, target, build_from, build_to):
        target_start = time.time()
        print(f"{idx}/{length} {target}")

        substitution = {
            'package': target,
            'binDir': build_to,
            'timestamp': datetime.now().strftime("%Y/%m/%d %H:%M:%S")
        }

        print("\tBuilding")
        os.system(f"GOOS=linux GOARCH=arm go build -o ./{build_to}{target} ./tmp/{build_from}{target}/*.go")

        print("\tGenerating dockerfile")
        with open('Dockerfile.template', 'r') as f:
            template = Template(f.read())
            result = template.substitute(substitution)

            with open(f"{build_to}/Dockerfile.{target}", 'w') as out:
                out.write(result)

        print(f"\tDone in {round((time.time() - target_start), 2)}s")

    print("Running build")

    # Copy the gomods directory into tmp
    print("Copying gomods")
    shutil.copytree("../gomods/", "./tmp/")

    # Build packages
    print("Building packages")
    for i, package in enumerate(packages):
        build_target(i+1, len(packages), package, "", "packages/")

    # Build balancers
    print("Building balancers")
    for i, balancer in enumerate(balancers):
        build_target(i+1, len(balancers), balancer, "balancers/", "balancers/")

    # Delete tmp directory
    print("Cleaning up")
    shutil.rmtree("./tmp/", ignore_errors=False, onerror=None)

elif args.target == "clean":
    print("Running clean")

    print("Removing 'packages'")
    shutil.rmtree("./packages/", ignore_errors=True, onerror=None)

    print("Removing 'balancers'")
    shutil.rmtree("./balancers/", ignore_errors=True, onerror=None)

elif args.target == "push":
    print("Pushing docker images")

    def push_image(idx, length, prefix, name, version, dockerFile):
        print(f"{idx}/{length} {prefix}{name}:{version}")
        print("\tBuilding")
        os.system(f"docker build --platform linux/arm64 -t mstk/{prefix}{name.lower()}:{version} -f {dockerFile} .")
        print("\tPushing")
        os.system(f"docker image push {prefix}{name.lower()}:{version}")

    for i, package in enumerate(packages):
        push_image(i+1, len(packages), "", package, "0.0.0", f"packages/Dockerfile.{package}")

    for i, balancer in enumerate(balancers):
        push_image(i+1, len(balancers), "balancers-", balancer, "0.0.0", f"balancers/Dockerfile.{balancer}")

print(f"Done in {round((time.time() - start_time), 2)}s")
