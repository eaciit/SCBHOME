package controllers

import (
	"eaciit/scbhome/webapp/web-main/models"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"io"
	"os"
	"path/filepath"
)

type DashboardController struct {
	*BaseController
}

func (c *DashboardController) Index(k *knot.WebContext) interface{} {
	c.SetupForHTML(k)

	return c.GetBaseData(k)
}

func (c *DashboardController) Master(k *knot.WebContext) interface{} {
	c.SetupForHTML(k)

	return c.GetBaseData(k)
}

func (c *DashboardController) GetData(k *knot.WebContext) interface{} {
	c.SetupForAJAX(k)

	res := c.GetBaseData(k)
	return res
}

func (c *DashboardController) GetPage(k *knot.WebContext) interface{} {
	c.SetupForAJAX(k)

	csr, err := c.Ctx.Connection.
		NewQuery().
		From(new(models.PageModel).TableName()).
		Select().
		Cursor(nil)
	if csr != nil {
		defer csr.Close()
	}
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	data := make([]models.PageModel, 0)
	err = csr.Fetch(&data, 0, true)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	return c.SetResultOK(data)
}

func (c *DashboardController) SavePage(k *knot.WebContext) interface{} {
	c.SetupForAJAX(k)

	payload := new(tk.M)
	fileHeader, formData, err := k.GetPayloadMultipart(payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	data := new(models.PageModel)
	tk.UnjsonFromString(formData["data"][0], data)

	if len(fileHeader["file"]) > 0 {
		file, handler, err := k.Request.FormFile("file")
		if err != nil {
			return c.SetResultError(err.Error(), nil)
		}

		defer file.Close()

		baseImagePath := Config.GetString("uploadpath")
		fileName := tk.RandomString(32) + filepath.Ext(handler.Filename)
		filePath := filepath.Join(baseImagePath, fileName)
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return c.SetResultError(err.Error(), nil)
		}

		defer f.Close()
		_, err = io.Copy(f, file)
		if err != nil {
			return c.SetResultError(err.Error(), nil)
		}

		data.Cover = fileName
	}

	if data.Id == "" {
		data.Id = tk.RandomString(32)
	}

	err = c.Ctx.Save(data)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	return c.SetResultOK(nil)
}

func (c *DashboardController) DeletePage(k *knot.WebContext) interface{} {
	c.SetupForAJAX(k)

	payload := struct {
		Ids []string
	}{}
	err := k.GetPayload(&payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	for _, eachId := range payload.Ids {
		eachModel := new(models.PageModel)
		eachModel.Id = eachId

		err = c.Ctx.Delete(eachModel)
		if err != nil {
			return c.SetResultError(err.Error(), nil)
		}
	}

	return c.SetResultOK(nil)
}

func (c *DashboardController) GetMasterPlatform(k *knot.WebContext) interface{} {
	c.SetupForAJAX(k)

	csr, err := c.Ctx.Connection.
		NewQuery().
		From(new(models.MasterPlatformModel).TableName()).
		Select().
		Cursor(nil)
	if csr != nil {
		defer csr.Close()
	}
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	data := make([]models.MasterPlatformModel, 0)
	err = csr.Fetch(&data, 0, true)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	return c.SetResultOK(data)
}

func (c *DashboardController) InsertPredefinedData(k *knot.WebContext) interface{} {
	c.SetupForAJAX(k)

	// platforms
	err := (func() error {
		platforms := []models.MasterPlatformModel{
			{Id: "P01", Name: "EACIIT App", Color: "#e67e22"},
			{Id: "P02", Name: "Exellerator", Color: "#03A9F4"},
			{Id: "P03", Name: "UAT", Color: "#27ae60"},
			{Id: "P04", Name: "Development", Color: "#d066e2"}}

		for _, each := range platforms {
			err := c.Ctx.Save(&each)
			if err != nil {
				return err
			}
		}

		return nil
	})()
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	// pages
	err = (func() error {
		const DESC = "Application for Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua"

		pages := []models.PageModel{
			{Id: "p1", ProjectName: "OCIR", Description: DESC, Cover: "michal-kubalczyk-257107.jpg",
				Platforms: []models.PagePlatform{
					{PlatformId: "P01", PlatformName: "EACIIT App", URL: "http://scbocir-dev.eaciit.com/live/dashboard/default", Username: "eaciit", Password: "Password.23"},
					{PlatformId: "P02", PlatformName: "Exellerator", URL: "https://ocir.exellerator.io/live/dashboard/default", Username: "007", Password: "Password.23"},
					{PlatformId: "P03", PlatformName: "UAT", URL: "", Username: "", Password: ""},
					{PlatformId: "P04", PlatformName: "Development", URL: "", Username: "", Password: ""},
				}},
			{Id: "p2", ProjectName: "IKEA", Description: DESC, Cover: "jeff-sheldon-3231.jpg",
				Platforms: []models.PagePlatform{
					{PlatformId: "P01", PlatformName: "EACIIT App", URL: "http://scbikea.eaciit.com/dashboard/default", Username: "eaciit", Password: "Password.1"},
					{PlatformId: "P02", PlatformName: "Exellerator", URL: "", Username: "", Password: ""},
					{PlatformId: "P03", PlatformName: "UAT", URL: "", Username: "", Password: ""},
					{PlatformId: "P04", PlatformName: "Development", URL: "", Username: "", Password: ""},
				}},
			{Id: "p3", ProjectName: "Super Connect 5", Description: DESC, Cover: "jason-ortego-5383.jpg",
				Platforms: []models.PagePlatform{
					{PlatformId: "P01", PlatformName: "EACIIT App", URL: "", Username: "", Password: ""},
					{PlatformId: "P02", PlatformName: "Exellerator", URL: "http://www.superconnect5.com/admin/dashboard/default", Username: "eaciit", Password: "Password.1"},
					{PlatformId: "P03", PlatformName: "UAT", URL: "", Username: "", Password: ""},
					{PlatformId: "P04", PlatformName: "Development", URL: "", Username: "", Password: ""},
				}},
			{Id: "p4", ProjectName: "BEF", Description: DESC, Cover: "andrew-neel-108081.jpg",
				Platforms: []models.PagePlatform{
					{PlatformId: "P01", PlatformName: "EACIIT App", URL: "http://scb-bef.eaciitapp.com/dashboard/default", Username: "", Password: ""},
					{PlatformId: "P02", PlatformName: "Exellerator", URL: "", Username: "", Password: ""},
					{PlatformId: "P03", PlatformName: "UAT", URL: "", Username: "", Password: ""},
					{PlatformId: "P04", PlatformName: "Development", URL: "", Username: "", Password: ""},
				}},
			{Id: "p5", ProjectName: "Sales Perf Dashboard", Description: DESC, Cover: "carlos-muza-84523.jpg",
				Platforms: []models.PagePlatform{
					{PlatformId: "P01", PlatformName: "EACIIT App", URL: "", Username: "", Password: ""},
					{PlatformId: "P02", PlatformName: "Exellerator", URL: "https://cbi.exellerator.io/dashboard/default", Username: "eaciit", Password: "Password.1"},
					{PlatformId: "P03", PlatformName: "UAT", URL: "", Username: "", Password: ""},
					{PlatformId: "P04", PlatformName: "Development", URL: "", Username: "", Password: ""},
				}},
			{Id: "p6", ProjectName: "Mobile Money", Description: DESC, Cover: "finance-mobile.jpg",
				Platforms: []models.PagePlatform{
					{PlatformId: "P01", PlatformName: "EACIIT App", URL: "http://scmm.eaciit.com/dashboard/default", Username: "eaciit", Password: "B1gD@ta"},
					{PlatformId: "P02", PlatformName: "Exellerator", URL: "", Username: "", Password: ""},
					{PlatformId: "P03", PlatformName: "UAT", URL: "", Username: "", Password: ""},
					{PlatformId: "P04", PlatformName: "Development", URL: "", Username: "", Password: ""},
				}},
			{Id: "p7", ProjectName: "TWIST", Description: DESC, Cover: "miguel-carraca-131178.jpg",
				Platforms: []models.PagePlatform{
					{PlatformId: "P01", PlatformName: "EACIIT App", URL: "", Username: "", Password: ""},
					{PlatformId: "P02", PlatformName: "Exellerator", URL: "https://twist.exellerator.io/", Username: "eaciit", Password: "Password.1"},
					{PlatformId: "P03", PlatformName: "UAT", URL: "", Username: "", Password: ""},
					{PlatformId: "P04", PlatformName: "Development", URL: "", Username: "", Password: ""},
				}},
			{Id: "p8", ProjectName: "FMI", Description: DESC, Cover: "Budgeting-and-Planning.jpg",
				Platforms: []models.PagePlatform{
					{PlatformId: "P01", PlatformName: "EACIIT App", URL: "", Username: "", Password: ""},
					{PlatformId: "P02", PlatformName: "Exellerator", URL: "https://fmi.exellerator.io", Username: "mid", Password: "Password.36"},
					{PlatformId: "P03", PlatformName: "UAT", URL: "", Username: "", Password: ""},
					{PlatformId: "P04", PlatformName: "Development", URL: "http://fmidev.eaciitapp.com", Username: "eaciit", Password: "knrkuvamdr@"},
				}},
			{Id: "p10", ProjectName: "RB BEF", Description: DESC, Cover: "calculator-calculation-insurance-finance-53621.jpeg",
				Platforms: []models.PagePlatform{
					{PlatformId: "P01", PlatformName: "EACIIT App", URL: "http://scb-rb.eaciitapp.com/web-main/login/default", Username: "rbuser", Password: "Password.1"},
					{PlatformId: "P02", PlatformName: "Exellerator", URL: "https://rbbizforum.exellerator.io", Username: "rbuser", Password: "Password.1"},
					{PlatformId: "P03", PlatformName: "UAT", URL: "", Username: "", Password: ""},
					{PlatformId: "P04", PlatformName: "Development", URL: "", Username: "", Password: ""},
				}},
			{Id: "p11", ProjectName: "FOT", Description: DESC, Cover: "calculator-calculation-insurance-finance-53621.jpeg",
				Platforms: []models.PagePlatform{
					{PlatformId: "P01", PlatformName: "EACIIT App", URL: "", Username: "", Password: ""},
					{PlatformId: "P02", PlatformName: "Exellerator", URL: "https://rbbizforum.exellerator.io", Username: "fotuser", Password: "Password.1"},
					{PlatformId: "P03", PlatformName: "UAT", URL: "", Username: "", Password: ""},
					{PlatformId: "P04", PlatformName: "Development", URL: "", Username: "", Password: ""},
				}},
			{Id: "p12", ProjectName: "COST MI", Description: DESC, Cover: "Budgeting-and-Planning.jpg",
				Platforms: []models.PagePlatform{
					{PlatformId: "P01", PlatformName: "EACIIT App", URL: "", Username: "", Password: ""},
					{PlatformId: "P02", PlatformName: "Exellerator", URL: "", Username: "", Password: ""},
					{PlatformId: "P03", PlatformName: "UAT", URL: "", Username: "", Password: ""},
					{PlatformId: "P04", PlatformName: "Development", URL: "http://fmidev.eaciitapp.com/costmi/login/default", Username: "eaciit", Password: "Password.1"},
				}},
			{Id: "p13", ProjectName: "TB Sales", Description: DESC, Cover: "Budgeting-and-Planning.jpg",
				Platforms: []models.PagePlatform{
					{PlatformId: "P01", PlatformName: "EACIIT App", URL: "", Username: "", Password: ""},
					{PlatformId: "P02", PlatformName: "Exellerator", URL: "https://tbsales.exellerator.io", Username: "eaciit", Password: "Password.1"},
					{PlatformId: "P03", PlatformName: "UAT", URL: "", Username: "", Password: ""},
					{PlatformId: "P04", PlatformName: "Development", URL: "", Username: "", Password: ""},
				}},
			{Id: "p14", ProjectName: "S2BPay", Description: DESC, Cover: "Budgeting-and-Planning.jpg",
				Platforms: []models.PagePlatform{
					{PlatformId: "P01", PlatformName: "EACIIT App", URL: "http://s2bpay.eaciitapp.com/login/default", Username: "eaciit", Password: "Password.1"},
					{PlatformId: "P02", PlatformName: "Exellerator", URL: "", Username: "", Password: ""},
					{PlatformId: "P03", PlatformName: "UAT", URL: "", Username: "", Password: ""},
					{PlatformId: "P04", PlatformName: "Development", URL: "", Username: "", Password: ""},
				}}}

		for _, each := range pages {
			err := c.Ctx.Save(&each)
			if err != nil {
				return err
			}
		}

		return nil
	})()
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	return c.SetResultOK(nil)
}
