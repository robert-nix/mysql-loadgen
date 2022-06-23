#!/usr/bin/env python3
import argparse
import matplotlib.pyplot as plt
import pandas as pd
import seaborn as sns


parser = argparse.ArgumentParser()
parser.add_argument('file', type=str)

args = parser.parse_args()

df = pd.read_csv(args.file, sep='\t')

df["t"] = df["runtimeMS"]/1000.0
df["QPS"] = df["queries"]*1000000.0/df["ittimeUS"]
df["Memory (MB)"] = df["rssPages"]*4.096/1024.0
df["CPU %"] = 10.0*df["utimeTicks"]
df["Table Opens (Hz)"] = 10.0*df["openedTables"]

for value in ("QPS", "Memory (MB)", "CPU %", "Table Opens (Hz)"):
    group = df.pivot("t", "table-definition-cache", value)
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
    plt.savefig(f"table-definition-cache ({value}).png", dpi=150)
