package realsearch

import "fmt"

func (c Client) GetScheduleV2(ruleID uint) ([]ScheduleV2, error) {
	req, err := c.NewRequest("GET", "schedule/v2/", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("rule_id", fmt.Sprint(ruleID))
	req.URL.RawQuery = q.Encode()

	var data []ScheduleV2
	return data, c.Do(req, &data)
}

func (c Client) GetTimeMachineRules() ([]TimeMachineItem, error) {
	req, err := c.NewRequest("GET", "time_machine/rules", nil)
	if err != nil {
		return nil, err
	}
	var data []TimeMachineItem
	return data, c.Do(req, &data)
}

func (c Client) GetTimeMachineItemRaws(itemID uint) ([]RawInfo, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("time_machine/item/%d/raws", itemID), nil)
	if err != nil {
		return nil, err
	}
	var data []RawInfo
	return data, c.Do(req, &data)
}
