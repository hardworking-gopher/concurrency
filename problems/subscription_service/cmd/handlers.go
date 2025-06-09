package main

import (
	"fmt"
	"github.com/pandaemoniumplaza/concurrency/problems/subscription_service/data"
	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	homePageTemplate     = "/home.page.gohtml"
	loginPageTemplate    = "/login.page.gohtml"
	registerPageTemplate = "/register.page.gohtml"
	plansPageTemplate    = "/plans.page.gohtml"
)

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, homePageTemplate, nil)
}

func (a *App) LoginPage(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, loginPageTemplate, nil)
}

func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	_ = a.Session.Destroy(r.Context())
	_ = a.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (a *App) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	_ = a.Session.RenewToken(r.Context())

	if err := r.ParseForm(); err != nil {
		a.ErrorLog.Println("failed to parse form")
		a.Session.Put(r.Context(), "error", "Something went wrong")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var (
		email = r.Form.Get("email")
		pwd   = r.Form.Get("password")
	)

	user, err := a.Models.User.GetByEmail(email)
	if err != nil {
		a.ErrorLog.Println("failed to query user", err)
		a.Session.Put(r.Context(), "error", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// TODO: Check whether user is activated

	b, err := user.PasswordMatches(pwd)
	if err != nil {
		a.ErrorLog.Println("failed to verify password", err)
		a.Session.Put(r.Context(), "error", "Something went wrong")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !b {
		a.ErrorLog.Println("password doesn't match")

		msg := Message{
			To:      email,
			Subject: "Login attempt",
			Data:    "Failed login attempt - invalid password",
		}

		a.sendEmail(msg)

		a.Session.Put(r.Context(), "error", "Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	a.Session.Put(r.Context(), "user", user)
	a.Session.Put(r.Context(), "userID", user.ID)

	a.Session.Put(r.Context(), "flash", "Successful login!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) RegisterPage(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, registerPageTemplate, nil)
}

func (a *App) PostRegisterPage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		a.ErrorLog.Println("failed to parse form")
		a.Session.Put(r.Context(), "error", "Something went wrong")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// TODO: Check whether user is already registered

	u := data.User{
		Email:     r.Form.Get("email"),
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Password:  r.Form.Get("password"), // TODO: Generate from hash?
		Active:    0,
		IsAdmin:   0,
	}

	if _, err := u.Insert(u); err != nil {
		a.ErrorLog.Println("failed to create user")
		a.Session.Put(r.Context(), "error", "Something went wrong")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	a.InfoLog.Println("user has been save to db")

	signedUrl := GenerateTokenFromString(fmt.Sprintf("http://localhost:8080/activate?email=%s", u.Email))

	a.InfoLog.Println("signed url", signedUrl)

	a.InfoLog.Println("sending activation email to user")

	msg := Message{
		To:       u.Email,
		Subject:  "Active your account",
		Template: "confirmation-email",
		Data:     template.HTML(signedUrl),
	}

	a.sendEmail(msg)

	a.Session.Put(r.Context(), "flash", "Confirmation email sent! Check your email.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (a *App) ActiveAccount(w http.ResponseWriter, r *http.Request) {
	a.InfoLog.Println("received activation request")

	url := r.RequestURI
	testURL := fmt.Sprintf("http://localhost:8080%s", url)
	okay := VerifyToken(testURL)

	if !okay {
		a.ErrorLog.Println("failed verify token from activation link")
		a.Session.Put(r.Context(), "error", "Invalid token")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	a.InfoLog.Println("url verification successful")

	u, err := a.Models.User.GetByEmail(r.URL.Query().Get("email"))
	if err != nil {
		a.ErrorLog.Println("failed to find user by email")
		a.Session.Put(r.Context(), "error", "No user found")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	u.Active = 1
	if err = u.Update(); err != nil {
		a.ErrorLog.Println("failed to change user status")

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	a.InfoLog.Println("user has been activated")

	a.Session.Put(r.Context(), "flash", "Account has been activated successfully.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) ChooseSubscription(w http.ResponseWriter, r *http.Request) {
	plans, err := a.Models.Plan.GetAll()
	if err != nil {
		a.ErrorLog.Println("failed to query plans")
		a.Session.Put(r.Context(), "error", "Something went wrong")
		return
	}

	dataMap := make(map[string]any)
	dataMap["plans"] = plans

	a.render(w, r, plansPageTemplate, &TemplateData{Data: dataMap})
}

func (a *App) Subscribe(w http.ResponseWriter, r *http.Request) {
	planId, _ := strconv.Atoi(r.URL.Query().Get("id"))

	plan, err := a.Models.Plan.GetOne(planId)
	if err != nil {
		a.ErrorLog.Println("failed to query plan")
		a.Session.Put(r.Context(), "error", "Something went wrong")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}

	user, ok := a.Session.Get(r.Context(), "user").(data.User)
	if !ok {
		a.ErrorLog.Println("failed to extract user from a session")
		a.Session.Put(r.Context(), "error", "Something went wrong")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}

	a.Wait.Add(1)

	go func() {
		defer a.Wait.Done()

		invoice, err := a.getInvoice(user, plan)
		if err != nil {
			a.ErrorChan <- err
			return
		}

		msg := Message{
			To:       user.Email,
			Subject:  "Your invoice",
			Data:     invoice,
			Template: "invoice",
		}

		a.sendEmail(msg)
	}()

	a.Wait.Add(1)

	go func() {
		defer a.Wait.Done()

		wd, _ := os.Getwd()

		filePath := fmt.Sprintf("%s/subscription_service/pdf/temp/%d_manual.pdf", wd, user.ID)

		err = a.generateManual(user, plan).OutputFileAndClose(filePath)
		if err != nil {
			a.ErrorChan <- err
			return
		}

		msg := Message{
			To:      user.Email,
			Subject: "Your manual",
			Data:    "Your manual is attached",
			AttachmentMap: map[string]string{
				"Manual.pdf": filePath,
			},
		}

		a.sendEmail(msg)
	}()

	if err = a.Models.Plan.SubscribeUserToPlan(user, *plan); err != nil {
		a.Session.Put(r.Context(), "error", "Error subscribing to plan")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}

	u, err := a.Models.User.GetOne(user.ID)
	if err != nil {
		a.Session.Put(r.Context(), "error", "Something went wrong")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}

	a.Session.Put(r.Context(), "user", u)
	a.Session.Put(r.Context(), "flash", "Subscribed!")
	http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
}

func (a *App) getInvoice(u data.User, plan *data.Plan) (string, error) {
	// dummy
	return plan.PlanAmountFormatted, nil
}

func (a *App) generateManual(u data.User, plan *data.Plan) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetMargins(10, 13, 10)

	importer := gofpdi.NewImporter()

	time.Sleep(5 * time.Second)

	wd, _ := os.Getwd()

	t := importer.ImportPage(pdf, fmt.Sprintf("%s/subscription_service/pdf/manual.pdf", wd), 1, "/MediaBox")
	pdf.AddPage()

	importer.UseImportedTemplate(pdf, t, 0, 0, 215.9, 0)

	pdf.SetX(75)
	pdf.SetY(150)

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 4, fmt.Sprintf("%s %s", u.FirstName, u.LastName), "", "C", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 4, fmt.Sprintf("%s User Guide", plan.PlanName), "", "C", false)

	return pdf

}
