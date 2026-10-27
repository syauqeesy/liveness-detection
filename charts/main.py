from pathlib import Path
import json

import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

# Letakkan script dan empat file JSON dalam folder yang sama.
DATA_DIR = Path(__file__).resolve().parent
OUTPUT_DIR = DATA_DIR / "latency_charts"
OUTPUT_DIR.mkdir(exist_ok=True)

FILES = {
    "server_inference": DATA_DIR / "normalized_logs.json",
    "client_inference": DATA_DIR / "normalized_console_logs.json",
    "managed_e2e": DATA_DIR / "normalized_managed_service_har_logs.json",
    "self_managed_e2e": DATA_DIR / "normalized_self_managed_har_logs.json",
}


def load_json(path):
    if not path.exists():
        raise FileNotFoundError(f"File tidak ditemukan: {path}")
    with path.open("r", encoding="utf-8") as f:
        return json.load(f)


def make_dataframe():
    # Inference latency dari log server: hanya event prediction completed.
    server_logs = load_json(FILES["server_inference"])
    server_rows = []
    for item in server_logs:
        if item.get("message") != "prediction completed":
            continue

        service = item.get("service", "")
        if service == "self-managed-prediction-service-liveness":
            approach = "Self-managed"
        elif service == "managed-prediction-service-liveness":
            approach = "Managed"
        else:
            continue

        value = item.get("duration_ms")
        if value is not None:
            server_rows.append({
                "metric": "Inference latency",
                "approach": approach,
                "latency_ms": float(value),
                "timestamp": item.get("timestamp"),
            })

    # Inference latency yang dicatat oleh browser/client.
    client_logs = load_json(FILES["client_inference"])
    for item in client_logs:
        value = item.get("duration_ms")
        if value is not None:
            server_rows.append({
                "metric": "Inference latency",
                "approach": "Client-side",
                "latency_ms": float(value),
                "timestamp": item.get("timestamp"),
            })

    # End-to-end latency dari HAR yang sudah dinormalisasi.
    for key, approach in [
        ("self_managed_e2e", "Self-managed"),
        ("managed_e2e", "Managed"),
    ]:
        logs = load_json(FILES[key])
        for item in logs:
            value = item.get("latency_ms")
            if value is not None:
                server_rows.append({
                    "metric": "End-to-end latency",
                    "approach": approach,
                    "latency_ms": float(value),
                    "timestamp": item.get("timestamp"),
                })

    df = pd.DataFrame(server_rows)
    df["timestamp"] = pd.to_datetime(df["timestamp"], errors="coerce", utc=True)
    df = df.dropna(subset=["latency_ms"])
    return df


def save_boxplot(df, metric, filename, order):
    subset = df[df["metric"] == metric].copy()
    if subset.empty:
        print(f"Tidak ada data untuk {metric}")
        return

    plt.figure(figsize=(9, 5.5))
    sns.boxplot(
        data=subset,
        x="approach",
        y="latency_ms",
        order=[x for x in order if x in subset["approach"].unique()],
        showfliers=True,
    )
    plt.title(f"{metric} by inference approach")
    plt.xlabel("Approach")
    plt.ylabel("Latency (ms)")
    plt.grid(axis="y", alpha=0.25)
    plt.tight_layout()
    plt.savefig(OUTPUT_DIR / filename, dpi=300, bbox_inches="tight")
    plt.close()


def save_mean_chart(df, metric, filename, order):
    subset = df[df["metric"] == metric].copy()
    if subset.empty:
        return

    summary = (
        subset.groupby("approach")["latency_ms"]
        .agg(mean="mean", std="std", count="count")
        .reindex([x for x in order if x in subset["approach"].unique()])
        .dropna(subset=["mean"])
    )

    plt.figure(figsize=(9, 5.5))
    plt.bar(
        summary.index,
        summary["mean"],
        yerr=summary["std"].fillna(0),
        capsize=5,
    )
    plt.title(f"Mean {metric} (error bars: ±1 standard deviation)")
    plt.xlabel("Approach")
    plt.ylabel("Mean latency (ms)")
    plt.grid(axis="y", alpha=0.25)
    plt.tight_layout()
    plt.savefig(OUTPUT_DIR / filename, dpi=300, bbox_inches="tight")
    plt.close()


def save_histogram(df, metric, filename, order):
    subset = df[df["metric"] == metric].copy()
    if subset.empty:
        return

    plt.figure(figsize=(9, 5.5))
    sns.histplot(
        data=subset,
        x="latency_ms",
        hue="approach",
        hue_order=[x for x in order if x in subset["approach"].unique()],
        bins="auto",
        element="step",
        stat="density",
        common_norm=False,
        fill=False,
    )
    plt.title(f"Distribution of {metric}")
    plt.xlabel("Latency (ms)")
    plt.ylabel("Density")
    plt.tight_layout()
    plt.savefig(OUTPUT_DIR / filename, dpi=300, bbox_inches="tight")
    plt.close()


def save_iteration_chart(df, metric, filename, order):
    subset = df[df["metric"] == metric].dropna(subset=["timestamp"]).copy()
    if subset.empty:
        return

    # Timestamp hanya dipakai sebagai urutan observasi dalam tiap pendekatan.
    subset = subset.sort_values("timestamp")
    subset["iteration"] = subset.groupby("approach").cumcount() + 1

    plt.figure(figsize=(10, 5.5))
    sns.lineplot(
        data=subset,
        x="iteration",
        y="latency_ms",
        hue="approach",
        hue_order=[x for x in order if x in subset["approach"].unique()],
        marker="o",
        linewidth=1,
        markersize=3,
    )
    plt.title(f"{metric} by measurement order")
    plt.xlabel("Iteration (ordered by timestamp within each approach)")
    plt.ylabel("Latency (ms)")
    plt.tight_layout()
    plt.savefig(OUTPUT_DIR / filename, dpi=300, bbox_inches="tight")
    plt.close()


def main():
    sns.set_theme(style="whitegrid", context="notebook")
    df = make_dataframe()

    inference_order = ["Self-managed", "Managed", "Client-side"]
    e2e_order = ["Self-managed", "Managed"]

    # Simpan data gabungan agar mudah diperiksa ulang.
    df.to_csv(OUTPUT_DIR / "latency_observations.csv", index=False)

    # Ringkasan statistik untuk laporan.
    summary = (
        df.groupby(["metric", "approach"])["latency_ms"]
        .agg(
            n="count",
            mean_ms="mean",
            median_ms="median",
            std_ms="std",
            min_ms="min",
            q1_ms=lambda s: s.quantile(0.25),
            q3_ms=lambda s: s.quantile(0.75),
            max_ms="max",
        )
        .reset_index()
    )
    summary.to_csv(OUTPUT_DIR / "latency_summary.csv", index=False)
    print("\nRingkasan statistik:")
    print(summary.round(2).to_string(index=False))

    save_boxplot(df, "Inference latency", "01_inference_latency_boxplot.png", inference_order)
    save_boxplot(df, "End-to-end latency", "02_end_to_end_latency_boxplot.png", e2e_order)

    save_mean_chart(df, "Inference latency", "03_inference_latency_mean.png", inference_order)
    save_mean_chart(df, "End-to-end latency", "04_end_to_end_latency_mean.png", e2e_order)

    save_histogram(df, "Inference latency", "05_inference_latency_histogram.png", inference_order)
    save_histogram(df, "End-to-end latency", "06_end_to_end_latency_histogram.png", e2e_order)

    save_iteration_chart(df, "Inference latency", "07_inference_latency_by_iteration.png", inference_order)
    save_iteration_chart(df, "End-to-end latency", "08_end_to_end_latency_by_iteration.png", e2e_order)

    print(f"\nSelesai. Grafik dan CSV tersimpan di: {OUTPUT_DIR}")


if __name__ == "__main__":
    main()
