package main

// All data structure will store in here
type GuiData struct {
	Image1_filepath string
	Image2_filepath string

	Processed1_filepath string
	Processed2_filepath string

	Compare_jpg_filepath string
	Compare_gif_filepath string

	Image1_size string
	Image2_size string
}

func NewGuiData(original1, processed1, filesize1, original2, processed2, filesize2, compare_jpg, compare_gif string) GuiData {
	return GuiData{
		Image1_filepath:      original1,
		Image2_filepath:      original2,
		Processed1_filepath:  processed1,
		Processed2_filepath:  processed2,
		Compare_jpg_filepath: compare_jpg,
		Compare_gif_filepath: compare_gif,
		Image1_size:          filesize1,
		Image2_size:          filesize2,
	}
}

type GuiDataSet struct {
	// An array of GuiData objects
	data []*GuiData
}

func (data *GuiDataSet) Get(index int) *GuiData {
	return data.data[index]
}

func (data *GuiDataSet) Add(item *GuiData) {
	data.data = append(data.data, item)
}

func (data *GuiDataSet) AddMultiple(items []*GuiData) {
	data.data = append(data.data, items...)
}

func (data *GuiDataSet) Remove(index int) *GuiData {
	temp := data.data[index]
	data.data = append(data.data[:index], data.data[index+1:]...)
	return temp
}

func NewGuiDataSetByFile(filepath string) *GuiDataSet {
	return &GuiDataSet{}
}
