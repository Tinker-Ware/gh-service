package usecases

import "github.com/gh-service/domain"

func (interactor GHInteractor) CreateFile(file domain.File, author domain.Author, username, repo string) error {

	err := interactor.GithubRepository.CreateFile(file, author, username, repo)
	if err != nil {
		return err
	}
	return nil
}

func (interactor GHInteractor) AddFiles(files []domain.File, author domain.Author, username, repo string) error {
	err := interactor.AddFiles(files, author, username, repo)
	if err != nil {
		return err
	}
	return nil
}
