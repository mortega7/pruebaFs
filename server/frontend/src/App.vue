<template>
  <nav class="navbar shadow navbar-expand-lg navbar-light bg-light rounded justify-content-between flex-row">
    <div class="container">
      <router-link class="navbar-brand float-left text-primary" to="/">FileServer Challenge</router-link>
      <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarToggler" aria-controls="navbarToggler" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbarToggler">
        <ul class="navbar-nav ms-auto mb-2 mb-lg-0">
          <li class="nav-item">
            <router-link class="btn btn-outline-primary" to="/"><i class="fa fa-home"></i> Go to index</router-link>
          </li>
        </ul>
      </div>
    </div>
  </nav>
  <router-view />
</template>
