#!/usr/bin/env python3
import argparse
import matplotlib.pyplot as plt
import pandas as pd
import seaborn as sns


parser = argparse.ArgumentParser()
parser.add_argument("file", type=str)
parser.add_argument("desc", type=str)

args = parser.parse_args()

df = pd.read_csv(args.file, sep="\t")

df["t"] = (df["runtimeMS"] / 1000.0).round(1)
df["QPS"] = df["queries"] * 1000000.0 / df["ittimeUS"]
df["Memory (MB)"] = df["rssPages"] * 4.096 / 1024.0
df["CPU %"] = 10.0 * df["utimeTicks"]

# log file has multiple series in order for each run
# create a run column to distinguish them by tracking when t resets
df["run"] = (df["t"].diff().fillna(0) < 0).cumsum()

for value in ("QPS", "Memory (MB)", "CPU %"):
    group = df.pivot(index="t", columns="run", values=value)
    print(group.head())
    sns.set_theme(style="darkgrid")
    g = sns.relplot(
        data=group,
        kind="line",
        dashes=False,
        aspect=3,
        lw=0.5,
    )
    g.set_axis_labels("Time (s)", value).tight_layout(w_pad=0)
    plt.savefig(f"{args.desc} ({value}).png", dpi=150)
