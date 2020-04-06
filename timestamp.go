package rtp

import "math"

type RtpTimestamp struct {
	//for pts
	ctime                        float64 //play time sec
	preRTPTimestamp              uint64
	rtpTimeIntervalMedian        uint64
	rtpTimeIntervalStatistics    map[uint64]int
	rtpTimeIntervalStatisticsCnt int
}

func NewRtpTimestamp() *RtpTimestamp {
	return &RtpTimestamp{
		0.0,
		0,
		0,
		make(map[uint64]int),
		0,
	}
}

//CalTimestampMs, rtpTimestamp为rtp头的时间戳; frequency为频率，视频一般为90000，返回值为从0开始的毫秒
func (rt *RtpTimestamp) CalTimestampMs(rtpTimestamp uint64, frequency float64) int64 {
	if rt.preRTPTimestamp == 0 { //start
	} else if uint64(rtpTimestamp) < rt.preRTPTimestamp { //reset
		rt.ctime = rt.ctime + float64(rt.rtpTimeIntervalMedian)/frequency
	} else if uint64(rtpTimestamp) > rt.preRTPTimestamp {
		interval := uint64(rtpTimestamp) - rt.preRTPTimestamp
		if cnt, ok := rt.rtpTimeIntervalStatistics[interval]; ok {
			rt.rtpTimeIntervalStatistics[interval] = cnt + 1
		} else {
			rt.rtpTimeIntervalStatistics[interval] = 1
		}
		rt.rtpTimeIntervalStatisticsCnt++
		if rt.rtpTimeIntervalStatisticsCnt > 20 {
			rt.rtpTimeIntervalStatisticsCnt = 0
			for k, v := range rt.rtpTimeIntervalStatistics {
				if _, ok := rt.rtpTimeIntervalStatistics[rt.rtpTimeIntervalMedian]; !ok {
					rt.rtpTimeIntervalMedian = k
				}
				if v > rt.rtpTimeIntervalStatistics[rt.rtpTimeIntervalMedian] {
					rt.rtpTimeIntervalMedian = k
				}
			}
			for k, _ := range rt.rtpTimeIntervalStatistics {
				delete(rt.rtpTimeIntervalStatistics, k)
			}
		}

		rt.ctime = rt.ctime + float64(interval)/frequency
	}
	rt.preRTPTimestamp = uint64(rtpTimestamp)

	return int64(uint64((27000000*rt.ctime)/300) % uint64(math.Pow(2, 33)))
}
