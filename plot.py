import seaborn as sns
import pandas as pd
import sys
import matplotlib.pyplot as plt

sns.set_theme()
ds = pd.read_table('plot.txt', names=['workers', 'rps_offered', 'ts', 'latency'])
plt.figure(figsize=(20, 14))
sns.displot(ds, x='latency', row='workers', col='rps_offered', binwidth=0.01)
plt.savefig('latency.png')

g = ds.groupby(['workers', 'rps_offered'])
rps_achieved = (g.ts.count() / g.ts.max()).to_frame('rps_achieved')

plt.figure(figsize=(10, 6))
sns.lineplot(rps_achieved, x='rps_offered', y='rps_achieved', hue='workers')
plt.savefig('rps_achieved.png')
