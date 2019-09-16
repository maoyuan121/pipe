package controller

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"github.com/b3log/gulu"
	"github.com/b3log/pipe/model"
	"github.com/b3log/pipe/service"
	"github.com/b3log/pipe/util"
	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

var states = map[string]string{}

// redirectGitHubLoginAction redirects to GitHub auth page.
func redirectGitHubLoginAction(c *gin.Context) {
	requestResult := gulu.Ret.NewResult()
	_, _, errs := gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Get(util.HacPaiURL+"/oauth/pipe/client2").
		Set("user-agent", model.UserAgent).Timeout(10 * time.Second).EndStruct(requestResult)
	if nil != errs {
		logger.Errorf("get oauth client id failed: %+v", errs)
		c.Status(http.StatusNotFound)

		return
	}
	if util.CodeOk != requestResult.Code {
		logger.Errorf("get oauth client id failed [code=%d, msg=%s]", requestResult.Code, requestResult.Msg)
		c.Status(http.StatusNotFound)

		return
	}
	data := requestResult.Data.(map[string]interface{})
	clientId := data["clientId"].(string)
	loginAuthURL := data["loginAuthURL"].(string)

	referer := c.Request.URL.Query().Get("referer")
	if "" == referer || !strings.Contains(referer, "://") {
		referer = model.Conf.Server + referer
	}
	if strings.HasSuffix(referer, "/") {
		referer = referer[:len(referer)-1]
	}
	state := gulu.Rand.String(16) + referer
	states[state] = state
	path := loginAuthURL + "?client_id=" + clientId + "&state=" + state + "&scope=public_repo,read:user,user:follow"

	logger.Infof("redirect to github [" + path + "]")

	c.Redirect(http.StatusSeeOther, path)
}

func githubCallbackAction(c *gin.Context) {
	logger.Infof("github callback [" + c.Request.URL.String() + "]")

	state := c.Query("state")
	if _, exist := states[state]; !exist {
		c.Status(http.StatusBadRequest)

		return
	}
	delete(states, state)

	referer := state[16:]
	if strings.Contains(referer, "__0") || strings.Contains(referer, "__1") {
		referer = referer[:len(referer)-len("__0")]
	}
	accessToken := c.Query("ak")
	githubUser := util.GitHubUserInfo(accessToken)
	if nil == githubUser {
		logger.Warnf("can not get user info with token [" + accessToken + "]")
		c.Status(http.StatusUnauthorized)

		return
	}

	githubId := githubUser["userId"].(string)
	userName := githubUser["userName"].(string)
	user := service.User.GetUserByGitHubId(githubId)
	if nil == user {
		if !service.Init.Inited() {
			user = &model.User{
				Name:      userName,
				AvatarURL: githubUser["userAvatarURL"].(string),
				B3Key:     githubId,
				GithubId:  githubId,
			}

			if err := service.Init.InitPlatform(user); nil != err {
				logger.Errorf("init platform via github login failed: " + err.Error())
				c.Status(http.StatusInternalServerError)

				return
			}
		} else {
			user = service.User.GetUserByName(userName)
			if nil == user {
				user = &model.User{
					Name:      userName,
					AvatarURL: githubUser["userAvatarURL"].(string),
					B3Key:     githubId,
					GithubId:  githubId,
				}

				if err := service.Init.InitBlog(user); nil != err {
					logger.Errorf("init blog via github login failed: " + err.Error())
					c.Status(http.StatusInternalServerError)

					return
				}
			} else {
				user.GithubId = githubId
				user.B3Key = githubId
				service.User.UpdateUser(user)
			}
		}
	}

	ownBlog := service.User.GetOwnBlog(user.ID)
	if nil == ownBlog {
		logger.Warnf("can not get user by name [" + userName + "]")
		c.Status(http.StatusNotFound)

		return
	}

	session := &util.SessionData{
		UID:     user.ID,
		UName:   user.Name,
		UB3Key:  user.B3Key,
		UAvatar: user.AvatarURL,
		URole:   ownBlog.UserRole,
		BID:     ownBlog.ID,
		BURL:    ownBlog.URL,
	}
	if err := session.Save(c); nil != err {
		logger.Errorf("saves session failed: " + err.Error())
		c.Status(http.StatusInternalServerError)
	}

	c.Redirect(http.StatusSeeOther, referer)
}
