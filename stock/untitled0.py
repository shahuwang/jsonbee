# -*- coding: utf-8 -*-
"""
Created on Tue Dec 13 00:26:06 2016

@author: Andrew
"""

import pandas as pd
import numpy as np

data = pd.read_csv('data_all.csv',index_col = 0)
data = data.sort_index()

#计算出收益率(用对数收益率方便计算)
return_data = (np.log(data)-np.log(data.shift(1)))


#得到第date的high low组合，data2应该是6个月的平均收益率的数据
def get_high_low_group(data2,date,k):
    temp = data2.loc[date].dropna()
    temp.sort_values(inplace = True)
    length = len(temp)
    low = list(temp[:int(0.1*length)].index)
    high =list(temp[-int(0.1*length):].index)
    
    return high,low

#得到K个月的组合的平均收益率
def get_return(data,group,date,k):
    return data[group].loc[date:].iloc[1:k+1].sum().mean()

#得到相应设置的收益，并保存为csv文件：
def get_final_results(data,k):
    dates = data.index
    #得到6个月的平均收益率
    data2 = pd.rolling_mean(data,window = k)
    temp1 = {}
    temp2 = {}

    for date in dates[k:len(dates)-k]:
        groups = get_high_low_group(data2,date,k)
        temp1[date] = (get_return(data,groups[0],date,k))
        temp2[date] = (get_return(data,groups[1],date,k))
    
    results = pd.DataFrame(pd.Series(temp1),columns = [['high']])
    results['low'] = pd.Series(temp2)
    results.to_csv('results_'+str(k)+'_month.csv')
    return results

#以6个月为周期
k = 6  
results = get_final_results(return_data,k)     
    
    
    
    
    


