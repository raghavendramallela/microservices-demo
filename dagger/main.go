// initial multi-stage docker build dagger module for java-gradle
package main

import (
	"context"
	"dagger/rgmodule/internal/dagger"
)

type Rgmodule struct{}

// build & publish image
func (rgmodule *Rgmodule) Build(
	ctx context.Context,
	// source code directory path containing gradle-wrapper(eg: --src=.)
	src *dagger.Directory,
) (string, error) {
	// build with gradle-wrapper
	builder := dag.Container().
		From("eclipse-temurin:21@sha256:b5fc642f67dbbd1c4ce811388801cb8480aaca8aa9e56fd6dcda362cfea113f1").
		WithWorkdir("/app").
		WithDirectory("/app", src).
		WithExec([]string{"chmod", "+x", "gradlew"}).
		WithExec([]string{"./gradlew", "downloadRepos"}).
		WithExec([]string{"./gradlew", "installDist"})

	// copy to jre
	jreImage := dag.Container().
		From("eclipse-temurin:21.0.4_7-jre-alpine@sha256:8cc1202a100e72f6e91bf05ab274b373a5def789ab6d9e3e293a61236662ac27").
		WithDirectory("/app", builder.Directory("/app")).
		WithExposedPort(9555).
		WithEntrypoint([]string{"/app/build/install/hipstershop/bin/AdService"})

	// publish image
	ttlshImage, err := jreImage.Publish(ctx, "ttl.sh/adservice:rg")
	if err != nil {
		return "", err
	}
	return ttlshImage, nil
}
