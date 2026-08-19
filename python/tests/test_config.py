import unittest

from eka_abdm.config import Config, Environment


class TestEnvironment(unittest.TestCase):
    def test_production_base_url(self):
        self.assertEqual(Environment.base_url(Environment.PRODUCTION), "https://api.eka.care")

    def test_development_base_url(self):
        self.assertEqual(
            Environment.base_url(Environment.DEVELOPMENT), "https://api.dev.eka.care"
        )

    def test_unknown_environment_defaults_to_production(self):
        self.assertEqual(Environment.base_url("staging"), "https://api.eka.care")


class TestConfigToken(unittest.TestCase):
    def test_static_token_default_empty(self):
        cfg = Config()
        self.assertEqual(cfg.token(), "")

    def test_static_token_is_returned(self):
        cfg = Config()
        cfg.set_authorization_token("abc123")
        self.assertEqual(cfg.token(), "abc123")

    def test_token_func_overrides_static_and_is_called_each_time(self):
        cfg = Config()
        cfg.set_authorization_token("static")
        calls = []

        def resolver():
            calls.append(1)
            return "dynamic-%d" % len(calls)

        cfg.set_token_func(resolver)
        self.assertEqual(cfg.token(), "dynamic-1")
        self.assertEqual(cfg.token(), "dynamic-2")  # resolved again, not cached here

    def test_default_base_url_derived_from_environment(self):
        cfg = Config(environment=Environment.DEVELOPMENT)
        self.assertEqual(cfg.base_url, "https://api.dev.eka.care")

    def test_explicit_base_url_overrides_environment(self):
        cfg = Config(environment=Environment.PRODUCTION, base_url="http://localhost:9999")
        self.assertEqual(cfg.base_url, "http://localhost:9999")


if __name__ == "__main__":
    unittest.main()
