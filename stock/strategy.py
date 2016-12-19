# -*- coding: utf-8 -*-
import pandas as pd
from datetime import datetime
from dateutil.relativedelta import relativedelta


EXCELS = [
        'TRD_Mnth 0.xls', 'TRD_Mnth1.xls',
        'TRD_Mnth2.xls', 'TRD_Mnth3.xls',
        'TRD_Mnth4.xls', 'TRD_Mnth5.xls'
        ]


def read_excel_to_df():
    df = pd.DataFrame()
    for f in EXCELS:
        mth = pd.read_excel(
                f, 'TRD_Mnth', skiprows=3,
                names=['code', 'date', 'price'],
                converters={'code': str, 'price': float})
        df = pd.concat([df, mth])
    df['date'] = pd.to_datetime(df['date'])
    df = df.sort(['code', 'date'])
    return df


def k_month_roi_before(df, year, month, interval=6):
    # year, month 为选股年月，如2000年7月
    # 计算该时间点过去六个月的回报率
    date = datetime(year, month, 1)
    lastm = date + relativedelta(months=-1)
    firstm = date + relativedelta(months=-(interval + 1))
    df1 = df.loc[df['date'].isin([firstm, lastm])]
    sub = df1['price'] - df1.shift(1)['price']
    sub = sub[1:]
    rate = sub / df.loc[df['date'] == firstm]['price']
    return pd.DataFrame({"code": pd.unique(df1['code']), "rate": rate})


def k_month_roi_after(df, codes, year, month, interval=6):
    df = df.loc[df['code'].isin(codes)]
    date = datetime(year, month, 1)
    lastm = date + relativedelta(months=1)
    firstm = date + relativedelta(months=interval - 1)
    df1 = df.loc[df['date'].isin([firstm, lastm])]
    sub = df1['price'] - df1.shift(1)['price']
    sub = sub[1:]
    rate = sub / df.loc[df['date'] == firstm]['price']
    return sum(rate) / len(rate)


def get_high_low_group(df, percent):
    # df 是回报率表
    count = int(len(df['code']) * percent)
    df = df.sort('rate')
    high = df.tail(count)
    low = df.head(count)
    return low['code'], high['code']


def run(year, month, interval, percent):
    df = read_excel_to_df()
    rate_df = k_month_roi_before(df, year, month, interval)
    low, high = get_high_low_group(rate_df, percent)
    print k_month_roi_after(df, low, year, month, interval)
    print k_month_roi_after(df, high, year, month, interval)


if __name__ == '__main__':
    run(2000, 7, 6, 0.1)
