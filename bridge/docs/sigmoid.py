import pandas as pd
import matplotlib.pyplot as plt

# Script to plot ALCB mint and burn data.

def printMint():
    df = pd.read_csv("bridge/mint1.csv", dtype="float")
    df["rate"] = df["mintedBToken"] / df["sentEth"]
    df2 = df[1:2200]
    df2.plot(kind="line", x="ether", y="rate", legend=None)
    plt.grid(True)
    plt.title("Mint Bonding Curve")
    plt.xlabel("Pool balance")
    plt.ylabel("ALCB/Ether conversion ratio")
    plt.show()


def printBurn():
    df = pd.read_csv("bridge/burn1.csv", dtype="float")
    df["rate"] = df["burnedBToken"] / df["receivedEth"]
    print(df)
    df2 = df[30:]
    df2.plot(kind="line", x="ether", y="rate", legend=None)
    plt.grid(True)
    plt.title("Burn Bonding Curve")
    plt.xlabel("Pool balance")
    plt.ylabel("ALCB/Ether conversion ratio")
    plt.show()


def printBurnAndMint():
    df = pd.read_csv("bridge/mint-burn1.csv", dtype="float")
    df["rate"] = df["sentEth"] / df["receivedEth"]
    print(df)
    df.plot(kind="line", x="ether", y="rate", legend=None)
    plt.grid(True)
    plt.title("Minting/Burning relation")
    plt.xlabel("Pool balance")
    plt.ylabel("Mint/Burn conversion ratio")
    plt.show()


def main():
    printMint()
    printBurn()
    printBurnAndMint()


if __name__ == "__main__":
    main()
