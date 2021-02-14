package main

import "testing"

func TestFixBudget(t *testing.T) {
	detail := `I&#039;m looking for a coder to develop an indicator on MT4 that makes moving averages alert when touched. But it also needs an extra parameter:<br /><br />
It needs to alert only when a certain amount of candles touch it. But the alert has to go off only when the next candle doesn&#039;t touch it at all. For example, say price hits the moving average and price retraced or goes through after 2 candles, I&#039;ll get an alert after the NEXT candle forms and closes having never touched the moving average. <br /><br />
The whole idea is that my strategy depends on the retracement off of a moving average once it&#039;s touched but if price loiters for too many candles, it&#039;s no longer valid. So I want to be alerted after price retraces and confirms retracement or go-through after the touch. I&#039;d like to be able to control the amount of candles to allow until the alert candle forms. So that if I set it up for a no more than 3 candle touch, then if a 4th consecutive candle forms and touches the moving average, no alert will be triggered. On the same token, if only 1 candle touched, the alert would trigger on the close of the 2nd candle that did not touch it.<br /><br />
I also need it to work with heiken ashi candle sticks as well as standard Japanese candle sticks.<br /><br />
Visual example attached<br /><br /><b>Budget</b>: $400
<br /><b>Posted On</b>: February 14, 2021 05:11 UTC<br /><b>Category</b>: Scripting &amp; Automation
<br /><b>Country</b>: United States
<br /><a href="https://www.upwork.com/jobs/MT4-indicator-coder_%7E0148a3ec5d2c3f6f72?source=rss">click to apply</a>`

	b, budget := checkEntryBudget(detail)

	if !b || budget != "400" {
		t.Error("Budget is not 20")
	}
}


func TestHourlyPrice(t *testing.T) {
	detail := `Hello<br /><br />
I would like to fork https://github.com/nwcd-samples/video-on-demand-on-aws&nbsp;&nbsp;and go in details of the whole thing and create poc for that , <br /><br />
this 5 days job if ur good with media on aws and lambda please apply<br /><br /><b>Hourly Range</b>: $13.00-$30.00

<br /><b>Posted On</b>: February 14, 2021 10:59 UTC<br /><b>Category</b>: Full Stack Development<br /><b>Skills</b>:Python,     AWS Lambda,     AWS Elemental    
<br /><b>Country</b>: United Arab Emirates
<br /><a href="https://www.upwork.com/jobs/Video-demand-aws-workflow_%7E015a7d03edd653d40b?source=rss">click to apply</a>`

	b, budget := getHourlyPrice(detail)

	if !b || budget != "400" {
		t.Error("Budget is not 20")
	}
}
